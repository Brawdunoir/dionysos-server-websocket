package requests

import (
	"encoding/json"
	"errors"
	"fmt"

	obj "github.com/Brawdunoir/goplay-server/objects"
	responses "github.com/Brawdunoir/goplay-server/responses"
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
		return responses.CreateResponse(nil, req.Code, err)
	}

	var request IRequest

	// Would be better to change r type and unmarshall/handle this at the end of switch
	switch req.Code {
	case NEW_CONNECTION:
		var r NewConnectionRequest
		err = json.Unmarshal(req.Payload, &r)
		request = r
	case NEW_ROOM:
		var r NewRoomRequest
		err = json.Unmarshal(req.Payload, &r)
		request = r
	case JOIN_ROOM:
		var r JoinRoomRequest
		err = json.Unmarshal(req.Payload, &r)
		request = r
	default:
		return responses.CreateResponse(nil, req.Code, errors.New("unknown code"))
	}
	if err != nil {
		return responses.CreateResponse(nil, req.Code, err)
	}

	err = request.Check()
	if err != nil {
		return responses.CreateResponse(nil, req.Code, err)
	}

	v, err := request.Handle(remoteAddr, conn, users, rooms)
	return responses.CreateResponse(v, req.Code, err)
}
