
// Tab tracking (for the icon)
var tabID = null;
var tabURL = "";
var tabPushtartNode = "";
var isInPushtartSetupTab = false;
function setAsActiveTab(url, pushtartNode){
  chrome.tabs.getSelected(null, function(tab) {
      tabID = tab.id;
      tabURL = url;
      tabPushtartNode = pushtartNode;
      isInPushtartSetupTab = true;
      chrome.browserAction.setIcon({"path": "pt_green.png"});
  });
}
chrome.tabs.onUpdated.addListener(function(tabId, changeInfo, tab) {
    if (tabId == tabID){
      chrome.browserAction.setIcon({"path": "pt_red.png"});
      isInPushtartSetupTab = false;
    }
});
chrome.tabs.onActivated.addListener(function(tabId, changeInfo, tab) {
  if (tabId != tabID){
    chrome.browserAction.setIcon({"path": "pt_red.png"});
    isInPushtartSetupTab = false;
  }
});






// Setting the headers etc
chrome.webRequest.onAuthRequired.addListener(onRequestIntercept, {urls: ["<all_urls>"]}, ["blocking"]);

function onRequestIntercept(details){
  console.log(details);
  if ((details.statusCode == 401) && (details.realm.indexOf("Pushtart:") != -1)){
    console.log("Pushtart proxy detected!", details.realm.substring(9));

    var url = new URL(details.url);
    var key = url.protocol+"!"+url.hostname+":"+url.port+"::"+details.realm.substring(9)
    if (allowedPushtartDomains[key]){
      return {authCredentials: {username: allowedPushtartDomains[key].username, password: allowedPushtartDomains[key].password}};
    }

    setAsActiveTab(details.url, details.realm.substring(9));
  }

}





//Message handler from the UI
chrome.runtime.onMessage.addListener(function(message)  {
    console.log("Got message:", message);
    if (message.msg == "setCredentials"){
      allowedPushtartDomains[message.domain.scheme+"!"+message.domain.host+":"+message.domain.port+"::"+message.domain.node] = {
        'domain': message.domain,
        'username': message.username,
        'password': message.password,
        'node': message.node,
      };
      saveDomains();
    }

    if (message.msg == "delete"){
      delete allowedPushtartDomains[message.key];
      saveDomains();
    }
});








var allowedPushtartDomains = {};

function saveDomains() {
  var settings = {
    'byDomain': allowedPushtartDomains,
  };
  chrome.storage.sync.set(settings, function() {
    console.log("Saved.")
  });
}

chrome.storage.onChanged.addListener(function(changes, namespace) {
  if (namespace == 'sync'){
    for (key in changes) {
      console.log("Data change to: " + key);
      var storageChange = changes[key];
      if (key == 'byDomain'){
        allowedPushtartDomains = storageChange.newValue;
      }
    }
  }
});

function loadStorage(){
  chrome.storage.sync.get(['byName', 'byDomain'], function(data){
    allowedPushtartDomains = data['byDomain'] || {};
    console.log("Storage loaded:", data);
  })
}

loadStorage();
