package main

import (
	"log"
	"net/http"

	"github.com/Brawdunoir/dionysos-server/objects"
	"github.com/Brawdunoir/dionysos-server/requests"
	"github.com/Brawdunoir/dionysos-server/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options
var users = objects.NewUsers()                                                             // list of users connected
var rooms = objects.NewRooms()                                                             // list of rooms registered
var slogger *zap.SugaredLogger

func main() {
	// Start the logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	slogger = logger.Sugar()

	// Start listening for websocket clients on port 8080
	slogger.Info("starting…")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", socketHandler)
	slogger.Info("listening…")
	http.ListenAndServe(":8080", router)
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Grab public IP of client
	publicAddr, err := utils.GetIPAdress(r)
	if err != nil {
		slogger.Errorw("cannot read a valid public ip from http header", "remoteAddr", r.RemoteAddr)
		return
	}

	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slogger.Errorw("cannot upgrade connection", "remoteAddr", r.RemoteAddr)
		return
	}
	defer conn.Close()

	slogger.Info("new peer connected")

	var req requests.Request

	for {
		err := conn.ReadJSON(&req)
		if err != nil {
			slogger.Errorw("cannot read json message", "remoteAddr", r.RemoteAddr)
			break
		}

		response, user := req.Handle(publicAddr, conn, users, rooms, slogger)

		// Send the response using mutex for concurrent calls to WriteJSON if the user
		// exists.
		if user != nil {
			user.SendJSON(response, slogger)
		} else {
			conn.WriteJSON(response)
		}

	}
}
