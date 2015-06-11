package main
 
import (
  "fmt"
  "github.com/gorilla/websocket"
  "html/template"
  "log"
  "net/http"
)

const (
  socketBufferSize  = 1024
  messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
  ReadBufferSize: socketBufferSize,
  WriteBufferSize: socketBufferSize,
  CheckOrigin: func(r *http.Request) bool { return true },
}

//==============================================================================
type echoBot struct {
  socket *websocket.Conn
}

func (bot *echoBot) run() {
  // Read the incoming message and send them back forever.
  for {
    _, msg, err := bot.socket.ReadMessage()
    if err != nil {
      return
    }
    received_message := string(msg)
    fmt.Printf("Received: %s\n", received_message)

    echo_message := fmt.Sprintf("echo of [%s]", msg)
    bot.socket.WriteMessage(websocket.TextMessage, []byte(echo_message))
    fmt.Printf("\tSent: %s\n\n", echo_message)
  }
}

// Handles request to the server: echoing back the input messages to the client.
func echoHandler(w http.ResponseWriter, r *http.Request) {
  socket, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }
  bot := &echoBot{socket: socket}
  go bot.run()
}

//==============================================================================
func clientHandler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("templates/client.html")
  t.Execute(w, r)
}
 
//==============================================================================
func main() {
  http.HandleFunc("/echo", echoHandler)
  http.HandleFunc("/", clientHandler)
  // FileServer handler for static files that need to be loaded (css, js).
  http.Handle(
    "/static/",
    http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic("Error: " + err.Error())
  }
}