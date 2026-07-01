(function () {
  "use strict";

  loadConfig().then(function (config) {
    var token = extractPairToken(location.pathname);
    if (!token) return;
    chrome.storage.local.set({
      provider_url: String(config.provider_url || "").replace(/\/+$/, ""),
      pair_token: token,
      pair_url: location.href,
      pair_token_updated_at: new Date().toISOString()
    });
    window.dispatchEvent(new CustomEvent("whatsmeow-passkey-token-saved", { detail: { token: token } }));
  }).catch(function () {});

  async function loadConfig() {
    var res = await fetch(chrome.runtime.getURL("config.json"), { cache: "no-store" });
    if (!res.ok) throw new Error("config.json indisponível");
    return res.json();
  }

  function extractPairToken(pathname) {
    var match = String(pathname || "").match(/^\/pair\/([^/?#]+)/);
    return match ? decodeURIComponent(match[1]) : "";
  }
})();
