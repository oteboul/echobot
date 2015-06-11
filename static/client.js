var socket = new WebSocket("ws://localhost:8080/echo");

$(document).ready(function(){
	socket.onmessage = function(e) {
	  $("p#echo").append(e.data + "<br>");
	}
  socket.onclose = function() {
    $("p#echo").append("<br><br>The connection has been close.<br><hr>");
  }

	$("input#sendButton").click(function() {
	  var message = $('textarea#message').val();
	  if (message) {
		  socket.send(message);
		  $('textarea#message').val("");	
	  }
	});
});