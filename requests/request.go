package requests

import (
	"encoding/json"
	"fmt"

	"github.com/Brawdunoir/dionysos-server/constants"
	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
	"gopkg.in/validator.v2"
)

type Request struct {
	Code    CodeType        `json:"code" validate:"len=3"`
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
func (r Request) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	err := r.Check()
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	var req IRequest

	// Would be better to use a map but it is kind of hard with current r.Handle method…
	switch r.Code {
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
	case KICK_PEER:
		req, err = createKickPeerRequest(r.Payload)
	case CHANGE_OWNER:
		req, err = createChangeOwnerRequest(r.Payload)
	default:
		response = res.NewErrorResponse(fmt.Sprint(constants.ERR_UNKNOW_CODE, ": ", r.Code), logger)
		return
	}
	if err != nil {
		response = res.NewErrorResponse(constants.ERR_INVALID_PAYLOAD, logger)
		return
	}

	err = validator.Validate(req)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	return req.Handle(client, users, rooms, logger)
}
