package requests

import (
	"encoding/json"
	"fmt"

	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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
func (req Request) Handle(publicAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) res.Response {
	err := req.Check()
	if err != nil {
		return res.NewErrorResponse(err.Error(), logger)
	}

	var request IRequest

	// Would be better to use a map but it is kind of hard with current r.Handle method…
	switch req.Code {
	case NEW_CONNECTION:
		request, err = createNewConnectionRequest(req.Payload)
	case NEW_ROOM:
		request, err = createNewRoomRequest(req.Payload)
	case JOIN_ROOM:
		request, err = createJoinRoomRequest(req.Payload)
	case JOIN_ROOM_ANSWER:
		request, err = createJoinRoomAnswerRequest(req.Payload)
	case NEW_MESSAGE:
		request, err = createNewMessageRequest(req.Payload)
	case CHANGE_USERNAME:
		request, err = createChangeUsernameRequest(req.Payload)
	default:
		return res.NewErrorResponse(fmt.Sprintf("unknown code: %s", req.Code), logger)
	}
	if err != nil {
		return res.NewErrorResponse("payload json not valid", logger)
	}

	err = request.Check()
	if err != nil {
		return res.NewErrorResponse(err.Error(), logger)
	}

	return request.Handle(publicAddr, conn, users, rooms, logger)
}
