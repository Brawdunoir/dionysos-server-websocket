package main

import (
	"log"
	"net/http"

	"github.com/Brawdunoir/dionysos-server/objects"
	"github.com/Brawdunoir/dionysos-server/requests"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {return true}} // use default options
var users = objects.NewUsers()      // list of users connected
var rooms = objects.NewRooms()      // list of rooms registered

func main() {
	log.Println("Starting…")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", socketHandler)
	log.Println("Listening…")
	http.ListenAndServe(":8080", router)
}


func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	log.Println(r.RemoteAddr, "Connected")

	for {
		var req requests.Request

		err := conn.ReadJSON(&req)
		if err != nil {
			log.Println("Error during JSON reading:", err)
			break
		}

		response := req.Handle(r.RemoteAddr, conn, users, rooms)
		if err != nil {
			log.Println("wrong request from client:", err)
		}
		conn.WriteJSON(response)
	}
}
