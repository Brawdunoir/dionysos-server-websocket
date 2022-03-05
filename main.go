package main

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Brawdunoir/dionysos-server/objects"
	"github.com/Brawdunoir/dionysos-server/requests"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options
var users = objects.NewUsers()                                                             // list of users connected
var rooms = objects.NewRooms()                                                             // list of rooms registered

func main() {
	log.Println("Starting…")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", socketHandler)
	log.Println("Listening…")
	http.ListenAndServe(":8080", router)
}

func getIPAdress(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	return ""
}

//ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}

var privateRanges = []ipRange{
	{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

// isPrivateSubnet - check to see if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	log.Println(getIPAdress(r))

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
