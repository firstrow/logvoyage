
var wsUri = "ws://"+window.location.host+":12345/ws"; 
var websocket = new WebSocket(wsUri); 
websocket.onopen = function(e) {console.log("Open")}; 
websocket.onclose = function(e) {console.log("Close")}; 
websocket.onmessage = function(e) {
    var obj = jQuery.parseJSON(e.data);
}; 
websocket.onerror = function(e) {console.log("error") };

function wsRegister(apiKey) {
    websocket.send(apiKey);
}
