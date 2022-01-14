package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

var users = map[string]*User{} // list of users connected
var rooms = map[string]*Room{} // list of rooms registered

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	log.Println("connected to:", r.RemoteAddr)

	for {
		var req Request

		err := conn.ReadJSON(&req)
		if err != nil {
			log.Println("Error during JSON reading:", err)
			break
		}

		_, err = req.handle(r.RemoteAddr, conn)
		if err != nil {
			log.Println("wrong request from client:", err)
			continue
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", socketHandler)
	http.ListenAndServe(":8080", router)
}
