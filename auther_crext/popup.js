function renderStatus(statusText) {
  document.getElementById('status').textContent = statusText;
}

document.addEventListener('DOMContentLoaded', function() {
 chrome.browserAction.setIcon({"path": "pt_green.png"});
 renderStatus('TODO load list of allowed stuff');
});

