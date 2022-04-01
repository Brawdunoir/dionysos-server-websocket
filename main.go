package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Brawdunoir/dionysos-server/objects"
	"github.com/Brawdunoir/dionysos-server/requests"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/Brawdunoir/dionysos-server/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// websocket upgrader
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// list of users connected
var users = objects.NewUsers()

// list of rooms registered
var rooms = objects.NewRooms()

// zap super logger
var slogger *zap.SugaredLogger

func main() {
	var err error

	// Load environment variables
	utils.LoadEnvironment()

	// Start the logger
	var logger *zap.Logger
	if os.Getenv(utils.KEY_ENVIRONMENT) == "PROD" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
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

	// Grab uuid generated and sent by client
	_, uuid, err := conn.ReadMessage()
	if err != nil {
		slogger.Error(err)
		return
	}

	// Add the client
	client := users.AddUser(publicAddr, string(uuid), conn, slogger)
	client.SendJSON(res.NewResponse(res.ConnectionResponse{UserID: client.ID}, slogger), slogger)

	var req requests.Request
	for {
		err := conn.ReadJSON(&req)
		if err != nil {
			slogger.Infow(err.Error(), "user", client.ID)
			requests.DisconnectPeer(client, users, rooms, slogger)
			break
		}

		response := req.Handle(client, users, rooms, slogger)

		// Send the response using mutex for concurrent calls to WriteJSON within Handlers.
		client.SendJSON(response, slogger)
	}
}
