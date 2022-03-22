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
func (r Request) Handle(publicAddr string, conn *websocket.Conn, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response, user *obj.User) {
	err := r.Check()
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	var req IRequest

	// Would be better to use a map but it is kind of hard with current r.Handle method…
	switch r.Code {
	case NEW_CONNECTION:
		req, err = createNewConnectionRequest(r.Payload)
	case NEW_ROOM:
		req, err = createNewRoomRequest(r.Payload)
	case JOIN_ROOM:
		req, err = createJoinRoomRequest(r.Payload)
	case JOIN_ROOM_ANSWER:
		req, err = createJoinRoomAnswerRequest(r.Payload)
	case NEW_MESSAGE:
		req, err = createNewMessageRequest(r.Payload)
	case CHANGE_USERNAME:
		req, err = createChangeUsernameRequest(r.Payload)
	case QUIT_ROOM:
		req, err = createQuitRoomRequest(r.Payload)
	case DISCONNECTION:
		req, err = createDisconnectionRequest(r.Payload)
	default:
		response = res.NewErrorResponse(fmt.Sprintf("unknown code: %s", r.Code), logger)
		return
	}
	if err != nil {
		response = res.NewErrorResponse("payload json not valid", logger)
		return
	}

	err = req.Check()
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	return req.Handle(publicAddr, conn, users, rooms, logger)
}
