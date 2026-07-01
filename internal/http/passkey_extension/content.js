(function () {
  "use strict";

  var config = null;
  var pending = {};
  var seq = 0;
  var lastRequestOptions = "";
  var currentToken = "";

  injectPageScript();
  init();

  function injectPageScript() {
    var script = document.createElement("script");
    script.src = chrome.runtime.getURL("inject.js");
    script.onload = function () { script.remove(); };
    (document.documentElement || document.head).appendChild(script);
  }

  async function init() {
    config = await loadConfig();
    var widget = createWidget();
    document.documentElement.appendChild(widget.root);
    chrome.runtime.onMessage.addListener(function (msg) {
      if (msg && msg.type === "wm-toggle") {
        widget.root.style.display = widget.root.style.display === "none" ? "block" : "none";
      }
    });
    widget.testButton.addEventListener("click", function () { void testCredential(widget); });
    widget.signButton.addEventListener("click", function () { void signPending(widget); });
    widget.saveButton.addEventListener("click", function () { void saveManualToken(widget); });
    void poll(widget);
    setInterval(function () { void poll(widget); }, 2500);
  }

  async function loadConfig() {
    var res = await fetch(chrome.runtime.getURL("config.json"), { cache: "no-store" });
    if (!res.ok) throw new Error("config.json indisponível");
    return res.json();
  }

  function createWidget() {
    var root = document.createElement("div");
    Object.assign(root.style, {
      position: "fixed",
      right: "16px",
      bottom: "16px",
      zIndex: 2147483647,
      width: "330px",
      borderRadius: "12px",
      overflow: "hidden",
      background: "#111b21",
      color: "#e9edef",
      boxShadow: "0 10px 30px rgba(0,0,0,.38)",
      fontFamily: "system-ui, -apple-system, BlinkMacSystemFont, Segoe UI, sans-serif",
      fontSize: "13px"
    });
    root.innerHTML =
      '<div style="display:flex;align-items:center;gap:8px;padding:10px 12px;background:#202c33;font-weight:700">' +
      '<span id="wm-dot" style="width:9px;height:9px;border-radius:999px;background:#8696a0"></span>' +
      '<span style="flex:1">Whatsmeow Passkey</span>' +
      '<button id="wm-close" type="button" title="Fechar" style="width:28px;height:28px;border:0;border-radius:6px;background:transparent;color:#e9edef;font-size:18px;cursor:pointer">×</button>' +
      '</div>' +
      '<div style="padding:12px">' +
      '<div id="wm-status" style="color:#aebac1;line-height:1.4;margin-bottom:10px">Inicializando...</div>' +
      '<input id="wm-token" type="text" placeholder="Cole o link /pair/... ou token" style="box-sizing:border-box;width:100%;height:36px;border:1px solid #3b4a54;border-radius:8px;background:#0b141a;color:#e9edef;padding:0 10px;margin-bottom:8px">' +
      '<button id="wm-save" type="button" style="width:100%;height:34px;border:1px solid #3b4a54;border-radius:8px;background:#172229;color:#e9edef;font-weight:700;cursor:pointer;margin-bottom:8px">Usar este link</button>' +
      '<button id="wm-sign" type="button" disabled style="width:100%;height:38px;border:0;border-radius:8px;background:#25d366;color:#052e1a;font-weight:800;cursor:pointer">Assinar desafio da conexão</button>' +
      '<button id="wm-test" type="button" style="width:100%;height:38px;border:1px solid #3b4a54;border-radius:8px;background:#172229;color:#e9edef;font-weight:700;cursor:pointer;margin-top:8px">Testar envio com minha passkey</button>' +
      '<div id="wm-result" style="margin-top:10px;color:#d1d7db;line-height:1.45;white-space:pre-wrap;word-break:break-word"></div>' +
      '</div>';
    root.querySelector("#wm-close").addEventListener("click", function () { root.style.display = "none"; });
    return {
      root: root,
      dot: root.querySelector("#wm-dot"),
      status: root.querySelector("#wm-status"),
      result: root.querySelector("#wm-result"),
      tokenInput: root.querySelector("#wm-token"),
      saveButton: root.querySelector("#wm-save"),
      signButton: root.querySelector("#wm-sign"),
      testButton: root.querySelector("#wm-test")
    };
  }

  async function poll(widget) {
    try {
      currentToken = await loadPairToken();
      if (document.activeElement !== widget.tokenInput) {
        widget.tokenInput.value = currentToken ? currentToken : "";
      }
      if (!currentToken) {
        lastRequestOptions = "";
        widget.dot.style.background = "#8696a0";
        widget.signButton.disabled = true;
        widget.status.textContent = "Abra o link /pair no provider ou cole o link/token aqui.";
        return;
      }
      var status = await providerJSON("/pair/" + encodeURIComponent(currentToken) + "/status");
      var options = status && status.passkey && status.passkey.request_options_json || "";
      lastRequestOptions = options;
      widget.dot.style.background = options ? "#25d366" : "#8696a0";
      widget.signButton.disabled = !options;
      if (status.link_status !== "active") {
        widget.status.textContent = "Link " + status.link_status + ". Gere um novo link no painel.";
      } else if (status.logged_in) {
        widget.status.textContent = "Conexão concluída.";
      } else if (options) {
        widget.status.textContent = "Desafio de passkey detectado. Clique para assinar.";
      } else {
        widget.status.textContent = "Extensão pronta. Aguardando o provider pedir passkey.";
      }
    } catch (err) {
      widget.dot.style.background = "#f15c6d";
      widget.signButton.disabled = true;
      widget.status.textContent = "Não consegui falar com o provider: " + err.message;
    }
  }

  async function signPending(widget) {
    if (!lastRequestOptions) {
      widget.result.textContent = "Ainda não há desafio de passkey para esta conexão.";
      return;
    }
    widget.signButton.disabled = true;
    widget.result.textContent = "Abrindo a chave de acesso do navegador...";
    var res = await runInPage("sign", { options: lastRequestOptions });
    if (!res.ok) {
      widget.result.textContent = "Falha na passkey: " + (res.name || "Error") + " " + (res.error || "");
      widget.signButton.disabled = false;
      return;
    }
    currentToken = await loadPairToken();
    if (!currentToken) {
      widget.result.textContent = "Abra o link /pair no provider ou cole o link/token antes de assinar.";
      widget.signButton.disabled = false;
      return;
    }
    var posted = await providerJSON("/pair/" + encodeURIComponent(currentToken) + "/passkey-prologue", {
      method: "POST",
      body: JSON.stringify({ webauthn_assertion: res.assertion })
    });
    widget.result.textContent = "Assertion enviada. Resposta: " + JSON.stringify(posted);
  }

  async function testCredential(widget) {
    widget.testButton.disabled = true;
    widget.result.textContent = "Abrindo sua carteira de senhas para um teste local...";
    try {
      var res = await runInPage("test", {});
      if (!res.ok) {
        widget.result.textContent = "Teste cancelado/falhou: " + (res.name || "Error") + " " + (res.error || "");
        return;
      }
      currentToken = await loadPairToken();
      if (!currentToken) {
        widget.result.textContent = "Abra o link /pair no provider ou cole o link/token antes de testar.";
        return;
      }
      var posted = await providerJSON("/pair/" + encodeURIComponent(currentToken) + "/passkey-test", {
        method: "POST",
        body: JSON.stringify({ webauthn_assertion: res.assertion })
      });
      widget.result.textContent = "Teste recebido pelo provider.\nrawId bytes: " + posted.raw_id_len + "\nassertion bytes: " + posted.assertion_len + "\nauthenticator: " + (res.authenticatorAttachment || "desconhecido");
    } catch (err) {
      widget.result.textContent = "Falha no teste: " + err.message;
    } finally {
      widget.testButton.disabled = false;
    }
  }

  async function saveManualToken(widget) {
    var token = extractPairToken(widget.tokenInput.value);
    if (!token) {
      widget.result.textContent = "Cole um link /pair válido ou o token.";
      return;
    }
    await chrome.storage.local.set({
      provider_url: String(config.provider_url || "").replace(/\/+$/, ""),
      pair_token: token,
      pair_token_updated_at: new Date().toISOString()
    });
    currentToken = token;
    widget.result.textContent = "Link salvo. Abra ou recarregue o WhatsApp Web para continuar.";
    await poll(widget);
  }

  async function loadPairToken() {
    var stored = await chrome.storage.local.get(["pair_token"]);
    return String(stored.pair_token || "").trim();
  }

  function extractPairToken(value) {
    var raw = String(value || "").trim();
    if (!raw) return "";
    try {
      var url = new URL(raw);
      var match = url.pathname.match(/^\/pair\/([^/?#]+)/);
      if (match) return decodeURIComponent(match[1]);
    } catch (err) {
      var directMatch = raw.match(/\/pair\/([^/?#]+)/);
      if (directMatch) return decodeURIComponent(directMatch[1]);
    }
    return raw;
  }

  function runInPage(action, payload) {
    var requestId = "wm-" + (++seq) + "-" + Date.now();
    return new Promise(function (resolve) {
      var timer = setTimeout(function () {
        if (pending[requestId]) {
          delete pending[requestId];
          resolve({ ok: false, name: "Timeout", error: "A chave de acesso não respondeu em tempo." });
        }
      }, 90000);
      pending[requestId] = function (data) {
        clearTimeout(timer);
        resolve(data);
      };
      window.postMessage(Object.assign({ __wmPasskey: true, action: action, requestId: requestId }, payload), "*");
    });
  }

  window.addEventListener("message", function (event) {
    var data = event.data;
    if (!data || data.__wmPasskey !== true || !data.requestId) return;
    if (data.action !== "sign-result" && data.action !== "test-result") return;
    var cb = pending[data.requestId];
    if (cb) {
      delete pending[data.requestId];
      cb(data);
    }
  });

  function providerJSON(path, opts) {
    var url = config.provider_url.replace(/\/+$/, "") + path;
    return chrome.runtime.sendMessage({
      type: "wm-fetch",
      method: opts && opts.method || "GET",
      url: url,
      body: opts && opts.body || undefined
    }).then(function (res) {
      if (!res || res.status === 0) throw new Error(res && res.error || "fetch failed");
      var text = res.text || "";
      var json = text ? JSON.parse(text) : {};
      if (!res.ok) throw new Error(json.error || json.message || text || ("HTTP " + res.status));
      return json;
    });
  }
})();
