package requests

import (
	"encoding/json"
	"errors"
	"fmt"

	obj "github.com/Brawdunoir/goplay-server/objects"
	"github.com/Brawdunoir/goplay-server/response"
	"github.com/gorilla/websocket"
)

type Request struct {
	Code    string          `json:"code"`
	Payload json.RawMessage `json:"payload"`
}

func (r Request) Check() error {
	var err error

	if r.Code == "" {
		err = fmt.Errorf("%w; code is empty", err)
	}
	if !json.Valid(r.Payload) {
		err = fmt.Errorf("%w; json not valid", err)
	}

	return err
}

// Handle creates a new request corresponding to the Code field
// and calls the Handle function on this new request
func (req Request) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) (interface{}, error) {
	err := req.Check()
	if err != nil {
		return response.CreateResponse(nil, err, req.Code)
	}

	var request IRequest

	// Would be better to change r type and unmarshall/handle this at the end of switch
	switch req.Code {
	case NEWCONNECTION:
		var r NewConnectionRequest
		err = json.Unmarshal(req.Payload, &r)
		request = r
	case NEWROOM:
		var r NewRoomRequest
		err = json.Unmarshal(req.Payload, &r)
		request = r
	case JOINROOM:
		var r JoinRoomRequest
		err = json.Unmarshal(req.Payload, &r)
		request = r
	default:
		return response.CreateResponse(nil, errors.New("unknown code"), req.Code)
	}
	if err != nil {
		return response.CreateResponse(nil, err, req.Code)
	}

	err = request.Check()
	if err != nil {
		return response.CreateResponse(nil, err, req.Code)
	}

	v, err := request.Handle(remoteAddr, conn, users, rooms)
	return response.CreateResponse(v, err, req.Code)
}
