(function () {
  "use strict";

  function b64url(buf) {
    var bytes = new Uint8Array(buf);
    var bin = "";
    for (var i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i]);
    return btoa(bin).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
  }

  function fromB64url(value) {
    var s = String(value || "").replace(/-/g, "+").replace(/_/g, "/");
    var pad = s.length % 4 ? "=".repeat(4 - (s.length % 4)) : "";
    var bin = atob(s + pad);
    var out = new Uint8Array(bin.length);
    for (var i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
    return out;
  }

  function serializeCredential(cred) {
    return {
      id: cred.id,
      rawId: b64url(cred.rawId),
      type: cred.type,
      response: {
        clientDataJSON: b64url(cred.response.clientDataJSON),
        authenticatorData: b64url(cred.response.authenticatorData),
        signature: b64url(cred.response.signature),
        userHandle: cred.response.userHandle ? b64url(cred.response.userHandle) : null
      }
    };
  }

  function normalizeOptions(rawOptions) {
    var opts = typeof rawOptions === "string" ? JSON.parse(rawOptions) : rawOptions;
    if (opts.challenge) opts.challenge = fromB64url(opts.challenge);
    if (Array.isArray(opts.allowCredentials)) {
      opts.allowCredentials.forEach(function (item) {
        if (item.id) item.id = fromB64url(item.id);
      });
    }
    return opts;
  }

  function randomChallenge() {
    var bytes = new Uint8Array(32);
    crypto.getRandomValues(bytes);
    return bytes;
  }

  window.addEventListener("message", function (event) {
    var data = event.data;
    if (!data || data.__wmPasskey !== true) return;

    if (data.action === "sign") {
      (async function () {
        try {
          var cred = await navigator.credentials.get({ publicKey: normalizeOptions(data.options) });
          window.postMessage({
            __wmPasskey: true,
            action: "sign-result",
            requestId: data.requestId,
            ok: true,
            assertion: serializeCredential(cred),
            authenticatorAttachment: cred.authenticatorAttachment || ""
          }, "*");
        } catch (err) {
          window.postMessage({
            __wmPasskey: true,
            action: "sign-result",
            requestId: data.requestId,
            ok: false,
            error: String(err && err.message ? err.message : err),
            name: String(err && err.name ? err.name : "Error")
          }, "*");
        }
      })();
    }

    if (data.action === "test") {
      (async function () {
        try {
          var cred = await navigator.credentials.get({
            publicKey: {
              challenge: randomChallenge(),
              rpId: "whatsapp.com",
              userVerification: "preferred",
              timeout: 60000
            }
          });
          window.postMessage({
            __wmPasskey: true,
            action: "test-result",
            requestId: data.requestId,
            ok: true,
            assertion: serializeCredential(cred),
            authenticatorAttachment: cred.authenticatorAttachment || ""
          }, "*");
        } catch (err) {
          window.postMessage({
            __wmPasskey: true,
            action: "test-result",
            requestId: data.requestId,
            ok: false,
            error: String(err && err.message ? err.message : err),
            name: String(err && err.name ? err.name : "Error")
          }, "*");
        }
      })();
    }
  });
})();
