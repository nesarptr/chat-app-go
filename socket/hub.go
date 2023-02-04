package socket

import (
	"fmt"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/nesarptr/chat-app-go/models"
)

var Clients = make(map[*websocket.Conn]*models.Client)
var Register = make(chan *websocket.Conn)
var UnRegister = make(chan *websocket.Conn)
var Broadcast = make(chan *models.Text)

func RunHub() {
	for {
		select {
		case c := <-Register:
			Clients[c].Mutex = new(sync.Mutex)
			fmt.Printf("connection registered for %s", Clients[c].UserName)
		case message := <-Broadcast:
			for c, user := range Clients {
				if user.UserName == message.SenderName || user.UserName == message.ReceiverName {
					go func(connection *websocket.Conn, u *models.Client) {
						u.Lock()
						defer u.Unlock()
						if u.IsClosing {
							return
						}
						if err := connection.WriteJSON(message); err != nil {
							u.IsClosing = true
							fmt.Println(err.Error())
							connection.WriteMessage(websocket.CloseMessage, []byte{})
							connection.Close()
							UnRegister <- connection
						}
					}(c, user)
				}
			}
		case c := <-UnRegister:
			username := Clients[c].UserName
			delete(Clients, c)
			fmt.Printf("connection unregistered for %s", username)
		}
	}
}
