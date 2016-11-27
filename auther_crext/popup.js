var domain = {};

function update(){
  var output = "<ul>";
  for (var key in chrome.extension.getBackgroundPage().allowedPushtartDomains) {
   if (chrome.extension.getBackgroundPage().allowedPushtartDomains.hasOwnProperty(key)) {
       var obj = chrome.extension.getBackgroundPage().allowedPushtartDomains[key];
       output += "<li>" + obj.domain.scheme + "//" + obj.domain.host + ":" + obj.domain.port + " (" + obj.domain.node + ")" +
                 "<a href=\"#\" id=\"" + encodeURIComponent(key) + "\">x</a></li>";
     }
   }
   output += "</ul>"
   document.getElementById('domainSection').innerHTML = output;

   var links = document.getElementsByTagName("a");
   for (var i = 0; i < links.length; i++) {
       (function () {
           var ln = links[i];
           var location = ln.id;
           ln.onclick = function () {
               chrome.runtime.sendMessage({
                 'msg': 'delete',
                 'key': decodeURIComponent(location),
               });
           };
       })();
   }

   if (chrome.extension.getBackgroundPage().isInPushtartSetupTab){
     document.getElementById('setup').style.display = 'block';
     var url = new URL(chrome.extension.getBackgroundPage().tabURL);
     document.getElementById('domain').textContent = url.protocol + "//" + url.hostname + ":" + url.port + " (" + chrome.extension.getBackgroundPage().tabPushtartNode + ")";
     domain = {
       'scheme': url.protocol,
       'host': url.hostname,
       'port': url.port,
       'node': chrome.extension.getBackgroundPage().tabPushtartNode,
     };
   } else {
     document.getElementById('setup').style.display = 'none';
   }
}


function setCredentialsPressed(){
  var username = document.getElementById('user').value;
  var password = document.getElementById('pas').value;
  chrome.runtime.sendMessage({
    'msg': 'setCredentials',
    'username': username,
    'password': password,
    'domain': domain,
  });
  document.getElementById('setup').style.display = 'none';
}

document.addEventListener('DOMContentLoaded', function() {
  update();
  document.getElementById('setCredentials').addEventListener('submit', setCredentialsPressed);
});


chrome.storage.onChanged.addListener(function(changes, namespace) {
  if (namespace == 'sync'){
    setTimeout(update, 100);
  }
});
