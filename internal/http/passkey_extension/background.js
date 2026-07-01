chrome.runtime.onMessage.addListener(function (msg, sender, sendResponse) {
  if (!msg || msg.type !== "wm-fetch") return;
  fetch(msg.url, {
    method: msg.method || "GET",
    headers: { "Content-Type": "application/json" },
    body: msg.body || undefined,
    cache: "no-store",
  })
    .then(function (res) {
      return res.text().then(function (text) {
        sendResponse({ status: res.status, ok: res.ok, text: text });
      });
    })
    .catch(function (err) {
      sendResponse({ status: 0, ok: false, error: String(err && err.message ? err.message : err) });
    });
  return true;
});

chrome.action.onClicked.addListener(function (tab) {
  if (tab && tab.id != null) {
    chrome.tabs.sendMessage(tab.id, { type: "wm-toggle" }, function () {
      void chrome.runtime.lastError;
    });
  }
});
