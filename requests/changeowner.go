package requests

import (
	"encoding/json"

	"github.com/Brawdunoir/dionysos-server/constants"
	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// ChangeOwnerRequest transfer ownership of room.
type ChangeOwnerRequest struct {
	NewOwnerID string `json:"newOwnerId" validate:"len=40"`
}

// Handle a ownership changement request from a room owner.
func (r ChangeOwnerRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	// Fetch room info
	newOwner, room, err := getUserByIdAndRoom(r.NewOwnerID, client.ID, users, rooms, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	// Assert user from this request is the legal owner of the room
	if room.OwnerID != client.ID {
		response = res.NewErrorResponse(constants.ERR_NO_PERM, logger)
		return
	}

	// Set new owner and notify peers
	err = room.SetOwnerID(newOwner, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}
	notifyPeers(rooms, room, logger)

	logger.Infow("change owner request success", "oldowner", client.ID, "oldownerusername", client.Name, "room", room.ID, "roomname", room.Name, "newowner", r.NewOwnerID)

	response = res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())}, logger)
	return
}

func (r ChangeOwnerRequest) Code() CodeType {
	return CHANGE_OWNER
}

func createChangeOwnerRequest(payload json.RawMessage) (r ChangeOwnerRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
