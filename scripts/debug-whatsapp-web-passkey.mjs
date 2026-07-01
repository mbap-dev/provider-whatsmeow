#!/usr/bin/env node

import fs from "node:fs";
import path from "node:path";

const marker = "__WA_PASSKEY_DEBUG__";

const args = parseArgs(process.argv.slice(2));
if (args.help) {
  console.log(`Usage:
  node scripts/debug-whatsapp-web-passkey.mjs [options]

Options:
  --host 127.0.0.1             Chrome DevTools host
  --port 9222                  Chrome DevTools port
  --url https://web.whatsapp.com/
  --out ./state/wa-passkey.ndjson
  --include-ws-payloads        Also save raw WebSocket frame payloads
  --dump-sources ./state/wa-sources
                               Save matching WhatsApp Web script snippets
  --search-sources passkey,webauthn
                               Comma-separated source search terms
  --dump-full-sources          Save full matching script files too
  --bridge-passkey             Fetch provider challenge, run WebAuthn in WhatsApp Web, submit assertion
  --provider-url http://localhost:8080
                               Provider base URL used with --bridge-passkey
  --session 5511999999999      Session ID used with --bridge-passkey
  --api-token wm_...           Optional bearer token for provider requests
  --api-version v1             API version for passkey-prologue submit
  --wait-passkey 120           Seconds to wait for provider passkey_required
  --help                       Show this help

Start Chrome first, for example on macOS:
  /Applications/Google\\ Chrome.app/Contents/MacOS/Google\\ Chrome \\
    --remote-debugging-port=9222 \\
    --user-data-dir=/tmp/wa-passkey-chrome \\
    --no-first-run \\
    https://web.whatsapp.com/
`);
  process.exit(0);
}
const host = args.host || "127.0.0.1";
const port = args.port || "9222";
const baseURL = `http://${host}:${port}`;
const targetURL = args.url || "https://web.whatsapp.com/";
const outPath = path.resolve(args.out || `./state/whatsapp-web-passkey-${Date.now()}.ndjson`);
const includeWSPayloads = Boolean(args["include-ws-payloads"]);
const sourceDumpDir = args["dump-sources"] ? path.resolve(args["dump-sources"]) : "";
const sourceSearchTerms = String(args["search-sources"] || "passkey,Passkey,credentials.get,navigator.credentials,publicKey,userVerification,webauthn,WebAuthn")
  .split(",")
  .map((term) => term.trim())
  .filter(Boolean);
const bridgePasskey = Boolean(args["bridge-passkey"]);
const providerURL = args["provider-url"] ? String(args["provider-url"]).replace(/\/+$/g, "") : "";
const bridgeSession = args.session ? String(args.session).trim() : "";
const apiToken = args["api-token"] || process.env.WHATSMEOW_API_TOKEN || "";
const apiVersion = args["api-version"] || "v1";
const waitPasskeySeconds = Number(args["wait-passkey"] || 0);

await fs.promises.mkdir(path.dirname(outPath), { recursive: true });
if (sourceDumpDir) {
  await fs.promises.mkdir(sourceDumpDir, { recursive: true });
}
const out = fs.createWriteStream(outPath, { flags: "a" });

process.on("SIGINT", () => {
  writeRecord({ kind: "process.interrupted" });
  out.end();
  process.exit(130);
});

writeRecord({
  kind: "capture.start",
  targetURL,
  chrome: baseURL,
  includeWSPayloads,
  sourceDumpDir: sourceDumpDir || undefined,
  sourceSearchTerms: sourceDumpDir ? sourceSearchTerms : undefined,
  bridgePasskey: bridgePasskey || undefined,
  providerURL: bridgePasskey ? providerURL : undefined,
  bridgeSession: bridgePasskey ? bridgeSession : undefined,
});

if (bridgePasskey && (!providerURL || !bridgeSession)) {
  throw new Error("--bridge-passkey requires --provider-url and --session");
}

const target = await findOrCreateTarget();
writeRecord({ kind: "target.selected", id: target.id, type: target.type, url: target.url });

const cdp = await connectCDP(target.webSocketDebuggerUrl);
const requestURLs = new Map();

cdp.onEvent = (event) => {
  const { method, params } = event;
  switch (method) {
    case "Runtime.consoleAPICalled":
      handleConsole(params);
      break;
    case "Network.webSocketCreated":
      requestURLs.set(params.requestId, params.url);
      writeRecord({
        kind: "cdp.websocket.created",
        requestId: params.requestId,
        url: params.url,
      });
      break;
    case "Network.webSocketClosed":
      writeRecord({
        kind: "cdp.websocket.closed",
        requestId: params.requestId,
        url: requestURLs.get(params.requestId),
        timestamp: params.timestamp,
      });
      break;
    case "Network.webSocketFrameError":
      writeRecord({
        kind: "cdp.websocket.error",
        requestId: params.requestId,
        url: requestURLs.get(params.requestId),
        timestamp: params.timestamp,
        errorMessage: params.errorMessage,
      });
      break;
    case "Network.webSocketFrameSent":
    case "Network.webSocketFrameReceived":
      writeWSFrame(method === "Network.webSocketFrameSent" ? "sent" : "received", params);
      break;
    case "Debugger.scriptParsed":
      if (sourceDumpDir) {
        void inspectScriptSource(params);
      }
      break;
  }
};

await cdp.send("Runtime.enable");
await cdp.send("Page.enable");
await cdp.send("Network.enable");
await cdp.send("Log.enable").catch(() => {});
if (sourceDumpDir) {
  await cdp.send("Debugger.enable");
}

const source = `(${installPasskeyDebugHook.toString()})(${JSON.stringify(marker)});`;
await cdp.send("Page.addScriptToEvaluateOnNewDocument", { source });
await cdp.send("Runtime.evaluate", {
  expression: source,
  awaitPromise: false,
  includeCommandLineAPI: false,
});

if (!target.url || !target.url.startsWith("https://web.whatsapp.com")) {
  await cdp.send("Page.navigate", { url: targetURL });
}

if (bridgePasskey) {
  await runPasskeyBridge();
  out.end();
  process.exit(0);
}

console.log(`Capturing WhatsApp Web passkey debug to ${outPath}`);
if (sourceDumpDir) {
  console.log(`Saving matching script snippets to ${sourceDumpDir}`);
}
console.log("Use the Chrome window opened with --remote-debugging-port, then run the WhatsApp Web passkey flow.");
console.log("Press Ctrl+C to stop.");

await new Promise(() => {});

async function runPasskeyBridge() {
  console.log(`Running passkey bridge for session ${bridgeSession}`);
  console.log(`Using WhatsApp Web target ${target.id}`);
  const status = await waitForPasskeyStatus();
  const requestOptionsJSON = getPasskeyRequestOptions(status);
  if (!requestOptionsJSON) {
    throw new Error(`provider status is passkey_required but request_options_json is empty: ${formatStatusSummary(status)}`);
  }
  writeRecord({
    kind: "bridge.challenge.loaded",
    session: bridgeSession,
    requestOptionsLength: requestOptionsJSON.length,
  });
  const result = await cdp.send("Runtime.evaluate", {
    expression: `(${runWebAuthnAssertionInPage.toString()})(${JSON.stringify(requestOptionsJSON)})`,
    awaitPromise: true,
    returnByValue: true,
    includeCommandLineAPI: false,
    userGesture: true,
  });
  if (result.exceptionDetails) {
    writeRecord({ kind: "bridge.webauthn.exception", exceptionDetails: result.exceptionDetails });
    throw new Error(result.exceptionDetails.text || "WebAuthn bridge failed");
  }
  const assertion = result.result?.value;
  if (!assertion?.id || !assertion?.response?.clientDataJSON) {
    writeRecord({ kind: "bridge.webauthn.invalid_result", assertion });
    throw new Error("WebAuthn did not return a valid assertion");
  }
  writeRecord({
    kind: "bridge.webauthn.assertion",
    session: bridgeSession,
    credentialID: assertion.id,
    hasUserHandle: assertion.response.userHandle != null,
  });
  const submit = await providerJSON(`/api/${encodeURIComponent(apiVersion)}/sessions/${encodeURIComponent(bridgeSession)}/passkey-prologue`, {
    method: "POST",
    body: JSON.stringify({
      credential_id: assertion.id,
      webauthn_assertion: assertion,
    }),
  });
  writeRecord({
    kind: "bridge.submitted",
    session: bridgeSession,
    response: submit,
  });
  console.log("Passkey assertion submitted to provider.");
}

async function waitForPasskeyStatus() {
  const deadline = Date.now() + waitPasskeySeconds * 1000;
  let lastStatus = null;
  while (true) {
    lastStatus = await providerJSON(`/sessions/${encodeURIComponent(bridgeSession)}/status`);
    const requestOptionsJSON = getPasskeyRequestOptions(lastStatus);
    if (lastStatus?.passkey?.required || requestOptionsJSON) {
      console.log(`Provider is waiting for passkey: ${formatStatusSummary(lastStatus)}`);
      return lastStatus;
    }
    if (!waitPasskeySeconds || Date.now() >= deadline) {
      writeRecord({
        kind: "bridge.no_pending_passkey",
        session: bridgeSession,
        status: summarizeStatus(lastStatus),
      });
      throw new Error(`provider status does not contain a pending passkey request: ${formatStatusSummary(lastStatus)}`);
    }
    console.log(`Waiting for provider passkey request: ${formatStatusSummary(lastStatus)}`);
    await sleep(2000);
  }
}

function getPasskeyRequestOptions(status) {
  return status?.passkey?.request_options_json || status?.passkey?.RequestOptionsJSON || "";
}

function summarizeStatus(status) {
  return {
    session_id: status?.session_id,
    connected: status?.connected,
    logged_in: status?.logged_in,
    device_jid: status?.device_jid,
    passkey_required: status?.passkey?.required || false,
    has_request_options: Boolean(getPasskeyRequestOptions(status)),
  };
}

function formatStatusSummary(status) {
  const summary = summarizeStatus(status);
  return `session=${summary.session_id || bridgeSession} connected=${summary.connected} logged_in=${summary.logged_in} passkey_required=${summary.passkey_required} has_request_options=${summary.has_request_options}`;
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function providerJSON(pathname, init = {}) {
  const headers = {
    Accept: "application/json",
    ...(init.body ? { "Content-Type": "application/json" } : {}),
    ...(apiToken ? { Authorization: `Bearer ${apiToken}` } : {}),
    ...(init.headers || {}),
  };
  const res = await fetch(`${providerURL}${pathname}`, {
    ...init,
    headers,
  });
  const text = await res.text();
  let payload = null;
  if (text) {
    try {
      payload = JSON.parse(text);
    } catch {
      payload = { raw: text };
    }
  }
  if (!res.ok) {
    throw new Error(`${pathname} returned ${res.status}: ${text}`);
  }
  return payload;
}

function writeWSFrame(direction, params) {
  const payload = params.response?.payloadData || "";
  writeRecord({
    kind: `cdp.websocket.frame.${direction}`,
    requestId: params.requestId,
    url: requestURLs.get(params.requestId),
    timestamp: params.timestamp,
    opcode: params.response?.opcode,
    mask: params.response?.mask,
    payloadLength: payload.length,
    payloadData: includeWSPayloads ? payload : undefined,
    payloadPreview: includeWSPayloads ? undefined : preview(payload),
  });
}

function handleConsole(params) {
  const values = (params.args || []).map((arg) => arg.value ?? arg.description ?? "");
  const markerIndex = values.indexOf(marker);
  if (markerIndex === -1) {
    return;
  }
  const raw = values[markerIndex + 1];
  if (typeof raw !== "string") {
    writeRecord({ kind: "injected.unparseable_console", values });
    return;
  }
  try {
    writeRecord({
      kind: "injected.passkey_debug",
      payload: JSON.parse(raw),
    });
  } catch {
    writeRecord({
      kind: "injected.passkey_debug.raw",
      payload: raw,
    });
  }
}

async function inspectScriptSource(params) {
  const url = params.url || "";
  if (!url.includes("web.whatsapp.com") && !url.includes("static.whatsapp.net")) {
    return;
  }
  if (!url.endsWith(".js") && !url.includes(".js?") && !url.includes("/lazy_loaded")) {
    return;
  }
  try {
    const result = await cdp.send("Debugger.getScriptSource", { scriptId: params.scriptId });
    const source = result.scriptSource || "";
    const matches = findSourceMatches(source, sourceSearchTerms);
    if (matches.length === 0) {
      return;
    }
    const fileBase = safeFileName(`${params.scriptId}-${url.split("/").pop() || "script.js"}`);
    const snippetsPath = path.join(sourceDumpDir, `${fileBase}.matches.txt`);
    const fullPath = path.join(sourceDumpDir, `${fileBase}.js`);
    await fs.promises.writeFile(snippetsPath, formatSourceMatches(url, matches), "utf8");
    if (args["dump-full-sources"]) {
      await fs.promises.writeFile(fullPath, source, "utf8");
    }
    writeRecord({
      kind: "source.matches",
      scriptId: params.scriptId,
      url,
      matchCount: matches.length,
      snippetsPath,
      fullPath: args["dump-full-sources"] ? fullPath : undefined,
      terms: [...new Set(matches.map((match) => match.term))],
    });
  } catch (error) {
    writeRecord({
      kind: "source.inspect_error",
      scriptId: params.scriptId,
      url,
      message: String(error?.message || error),
    });
  }
}

function findSourceMatches(source, terms) {
  const matches = [];
  for (const term of terms) {
    let fromIndex = 0;
    while (fromIndex < source.length) {
      const index = source.indexOf(term, fromIndex);
      if (index === -1) {
        break;
      }
      matches.push({
        term,
        index,
        snippet: source.slice(Math.max(0, index - 900), Math.min(source.length, index + term.length + 1400)),
      });
      fromIndex = index + term.length;
      if (matches.length >= 80) {
        return matches;
      }
    }
  }
  return matches.sort((a, b) => a.index - b.index);
}

function formatSourceMatches(url, matches) {
  return [
    `url: ${url}`,
    `matches: ${matches.length}`,
    "",
    ...matches.map((match, index) => [
      `===== match ${index + 1}: ${match.term} @ ${match.index} =====`,
      match.snippet,
      "",
    ].join("\n")),
  ].join("\n");
}

function safeFileName(value) {
  return value.replace(/[^a-zA-Z0-9._-]+/g, "_").slice(0, 180);
}

function writeRecord(record) {
  out.write(JSON.stringify({
    capturedAt: new Date().toISOString(),
    ...record,
  }) + "\n");
}

function preview(value) {
  if (typeof value !== "string" || value === "") {
    return "";
  }
  if (/^[\x20-\x7E]+$/.test(value)) {
    return value.slice(0, 240);
  }
  return `<non-printable length=${value.length}>`;
}

async function findOrCreateTarget() {
  const targets = await cdpHTTP("/json/list");
  const existing = targets.find((item) => item.type === "page" && item.url?.startsWith("https://web.whatsapp.com"));
  if (existing?.webSocketDebuggerUrl) {
    return existing;
  }
  const created = await cdpHTTP(`/json/new?${encodeURIComponent(targetURL)}`, { method: "PUT" }).catch(() => null);
  if (created?.webSocketDebuggerUrl) {
    return created;
  }
  const fallback = await cdpHTTP(`/json/new?${encodeURIComponent(targetURL)}`);
  if (!fallback?.webSocketDebuggerUrl) {
    throw new Error("Chrome did not return a page websocket debugger URL");
  }
  return fallback;
}

async function cdpHTTP(endpoint, init) {
  const res = await fetch(`${baseURL}${endpoint}`, init);
  if (!res.ok) {
    throw new Error(`${endpoint} returned ${res.status}: ${await res.text()}`);
  }
  return res.json();
}

async function connectCDP(wsURL) {
  const ws = new WebSocket(wsURL);
  const pending = new Map();
  let nextID = 1;
  const cdp = {
    onEvent: () => {},
    send(method, params = {}) {
      const id = nextID++;
      ws.send(JSON.stringify({ id, method, params }));
      return new Promise((resolve, reject) => {
        pending.set(id, { resolve, reject, method });
      });
    },
  };
  await new Promise((resolve, reject) => {
    ws.addEventListener("open", resolve, { once: true });
    ws.addEventListener("error", reject, { once: true });
  });
  ws.addEventListener("message", (event) => {
    const message = JSON.parse(event.data);
    if (message.id) {
      const item = pending.get(message.id);
      if (!item) {
        return;
      }
      pending.delete(message.id);
      if (message.error) {
        item.reject(new Error(`${item.method}: ${message.error.message}`));
      } else {
        item.resolve(message.result);
      }
      return;
    }
    cdp.onEvent(message);
  });
  return cdp;
}

function parseArgs(argv) {
  const parsed = {};
  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (!arg.startsWith("--")) {
      continue;
    }
    const key = arg.slice(2);
    if (key === "include-ws-payloads" || key === "dump-full-sources" || key === "bridge-passkey" || key === "help") {
      parsed[key] = true;
      continue;
    }
    parsed[key] = argv[i + 1];
    i++;
  }
  return parsed;
}

async function runWebAuthnAssertionInPage(requestOptionsJSON) {
  const options = JSON.parse(requestOptionsJSON);

  function fromBase64URL(value) {
    const base64 = String(value).replace(/-/g, "+").replace(/_/g, "/");
    const padded = base64 + "=".repeat((4 - base64.length % 4) % 4);
    const binary = atob(padded);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    return bytes.buffer;
  }

  function toBytes(value) {
    if (value instanceof ArrayBuffer) {
      return new Uint8Array(value);
    }
    if (ArrayBuffer.isView(value)) {
      return new Uint8Array(value.buffer, value.byteOffset, value.byteLength);
    }
    return null;
  }

  function toBase64URL(value) {
    const bytes = toBytes(value);
    if (!bytes) {
      return null;
    }
    let binary = "";
    for (const byte of bytes) {
      binary += String.fromCharCode(byte);
    }
    return btoa(binary).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/g, "");
  }

  if (typeof options.challenge === "string") {
    options.challenge = fromBase64URL(options.challenge);
  }
  if (Array.isArray(options.allowCredentials)) {
    options.allowCredentials = options.allowCredentials.map((credential) => ({
      ...credential,
      id: typeof credential.id === "string" ? fromBase64URL(credential.id) : credential.id,
    }));
  }

  const credential = await navigator.credentials.get({ publicKey: options });
  if (!credential) {
    throw new Error("navigator.credentials.get returned null");
  }
  const response = credential.response || {};
  return {
    id: credential.id,
    rawId: toBase64URL(credential.rawId),
    type: credential.type,
    response: {
      clientDataJSON: toBase64URL(response.clientDataJSON),
      authenticatorData: toBase64URL(response.authenticatorData),
      signature: toBase64URL(response.signature),
      userHandle: response.userHandle ? toBase64URL(response.userHandle) : null,
    },
  };
}

function installPasskeyDebugHook(markerValue) {
  if (window.__waPasskeyDebugInstalled) {
    return;
  }
  window.__waPasskeyDebugInstalled = true;

  const encoder = new TextEncoder();
  const decoder = new TextDecoder();
  let lastPasskeyResultAt = 0;

  function emit(payload) {
    try {
      console.info(markerValue, JSON.stringify({
        at: new Date().toISOString(),
        href: location.href,
        ...payload,
      }));
    } catch (error) {
      console.info(markerValue, JSON.stringify({
        kind: "hook.emit_error",
        message: String(error && error.message || error),
      }));
    }
  }

  function toBytes(value) {
    if (value instanceof ArrayBuffer) {
      return new Uint8Array(value);
    }
    if (ArrayBuffer.isView(value)) {
      return new Uint8Array(value.buffer, value.byteOffset, value.byteLength);
    }
    return null;
  }

  function toBase64URL(value) {
    const bytes = toBytes(value);
    if (!bytes) {
      return null;
    }
    let binary = "";
    for (const byte of bytes) {
      binary += String.fromCharCode(byte);
    }
    return btoa(binary).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/g, "");
  }

  function textFromBuffer(value) {
    const bytes = toBytes(value);
    if (!bytes) {
      return null;
    }
    try {
      return decoder.decode(bytes);
    } catch {
      return null;
    }
  }

  function serialize(value, depth = 0) {
    if (depth > 8) {
      return "[max-depth]";
    }
    if (value == null || typeof value === "string" || typeof value === "number" || typeof value === "boolean") {
      return value;
    }
    const bytes = toBytes(value);
    if (bytes) {
      return {
        type: value.constructor?.name || "ArrayBuffer",
        byteLength: bytes.byteLength,
        base64url: toBase64URL(bytes),
      };
    }
    if (Array.isArray(value)) {
      return value.map((item) => serialize(item, depth + 1));
    }
    if (typeof value === "object") {
      const out = {};
      for (const key of Object.keys(value)) {
        out[key] = serialize(value[key], depth + 1);
      }
      return out;
    }
    return String(value);
  }

  async function serializeCredential(credential) {
    if (!credential) {
      return null;
    }
    const response = credential.response || {};
    const clientDataJSONText = textFromBuffer(response.clientDataJSON);
    return {
      id: credential.id,
      type: credential.type,
      authenticatorAttachment: credential.authenticatorAttachment,
      rawId: toBase64URL(credential.rawId),
      response: {
        authenticatorData: toBase64URL(response.authenticatorData),
        clientDataJSON: toBase64URL(response.clientDataJSON),
        clientDataJSONText,
        signature: toBase64URL(response.signature),
        userHandle: response.userHandle ? toBase64URL(response.userHandle) : null,
      },
      clientExtensionResults: typeof credential.getClientExtensionResults === "function"
        ? credential.getClientExtensionResults()
        : null,
    };
  }

  const credentials = navigator.credentials;
  if (credentials?.get && !credentials.get.__waPasskeyDebugWrapped) {
    const originalGet = credentials.get.bind(credentials);
    const wrappedGet = async function(options) {
      const startedAt = Date.now();
      emit({
        kind: "credentials.get.call",
        options: serialize(options),
      });
      try {
        const credential = await originalGet(options);
        const serialized = await serializeCredential(credential);
        lastPasskeyResultAt = Date.now();
        window.__waPasskeyDebugLastCredential = serialized;
        emit({
          kind: "credentials.get.result",
          elapsedMs: Date.now() - startedAt,
          credential: serialized,
        });
        return credential;
      } catch (error) {
        emit({
          kind: "credentials.get.error",
          elapsedMs: Date.now() - startedAt,
          name: error?.name,
          message: error?.message,
          stack: error?.stack,
        });
        throw error;
      }
    };
    wrappedGet.__waPasskeyDebugWrapped = true;
    Object.defineProperty(credentials, "get", {
      value: wrappedGet,
      configurable: true,
      writable: true,
    });
  }

  if (window.WebSocket?.prototype?.send && !window.WebSocket.prototype.send.__waPasskeyDebugWrapped) {
    const originalSend = window.WebSocket.prototype.send;
    const wrappedSend = function(data) {
      const bytes = toBytes(data);
      const afterPasskeyMs = lastPasskeyResultAt ? Date.now() - lastPasskeyResultAt : null;
      emit({
        kind: "websocket.send",
        url: this.url,
        readyState: this.readyState,
        dataType: data?.constructor?.name || typeof data,
        byteLength: bytes ? bytes.byteLength : (typeof data === "string" ? encoder.encode(data).byteLength : null),
        textPreview: typeof data === "string" ? data.slice(0, 240) : null,
        afterPasskeyMs: afterPasskeyMs != null && afterPasskeyMs < 30000 ? afterPasskeyMs : null,
      });
      return originalSend.apply(this, arguments);
    };
    wrappedSend.__waPasskeyDebugWrapped = true;
    Object.defineProperty(window.WebSocket.prototype, "send", {
      value: wrappedSend,
      configurable: true,
      writable: true,
    });
  }

  emit({
    kind: "hook.installed",
    userAgent: navigator.userAgent,
  });
}
