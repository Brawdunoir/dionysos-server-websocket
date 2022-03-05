package requests

import (
	"encoding/json"
	"fmt"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
)

type Request struct {
	Code    CodeType        `json:"code"`
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
func (req Request) Handle(remoteAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms) res.Response {
	err := req.Check()
	if err != nil {
		return res.NewErrorResponse(err.Error())
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
	case JOIN_ROOM_ANSWER:
		var r JoinRoomAnswerRequest
		err = json.Unmarshal(req.Payload, &r)
		request = r
	default:
		return res.NewErrorResponse(fmt.Sprintf("unknown code: %s", req.Code))
	}
	if err != nil {
		return res.NewErrorResponse("payload json not valid")
	}

	err = request.Check()
	if err != nil {
		return res.NewErrorResponse(err.Error())
	}

	return request.Handle(remoteAddr, conn, users, rooms)
}
