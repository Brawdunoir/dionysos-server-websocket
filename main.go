package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

var users = map[string]User{} // list of users connected
var rooms = map[string]Room{} // list of rooms registered

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	for {
		// We’re reading keyword to know the action to perform next
		_, action, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Println("Received:", action, "from", r.RemoteAddr)

		// Perform desired action
		switch {
		case same(action, NEWCONNECTION):
			err := HandleConnection(r.RemoteAddr, conn)
			if err != nil {
				log.Println(err)
				// TODO send error to client
			}
		case same(action, NEWROOM):
			err := HandleNewRoom(r.RemoteAddr, conn)
			if err != nil {
				log.Println(err)
				// TODO send error to client
			}
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/socket", gin.WrapF(socketHandler))
	router.Run("localhost:8080")
}
