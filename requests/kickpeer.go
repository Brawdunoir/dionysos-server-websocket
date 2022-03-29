package requests

import (
	"encoding/json"
	"fmt"

	"github.com/Brawdunoir/dionysos-server/constants"
	obj "github.com/Brawdunoir/dionysos-server/objects"
	res "github.com/Brawdunoir/dionysos-server/responses"
	"go.uber.org/zap"
)

// KickPeerRequest allows the room's owner to remove a peer from his room.
type KickPeerRequest struct {
	PeerID string `json:"peerId"`
}

func (r KickPeerRequest) Check() error {
	var err error

	if r.PeerID == "" {
		err = fmt.Errorf("%w; peerId is empty", err)
	}

	return err
}

// Handles a kick peer in a room.
func (r KickPeerRequest) Handle(client *obj.User, users *obj.Users, rooms *obj.Rooms, logger *zap.SugaredLogger) (response res.Response) {
	// Fetch the peer to kick
	peer, room, err := getUserByIdAndRoom(r.PeerID, client.RoomID, users, rooms, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	// Assert user from this request is the legal owner of the room
	if room.OwnerID != client.ID {
		response = res.NewErrorResponse(constants.ERR_NO_PERM, logger)
		return
	}

	// Kick the peer from the room
	err = removeUserFromRoom(peer, users, rooms, logger)
	if err != nil {
		response = res.NewErrorResponse(err.Error(), logger)
		return
	}

	logger.Infow("kick peer request", "requester", client.ID, "requesterUsername", client.Name, "kickedPeer", peer.ID, "kickedPeerUsername", peer.Name, "room", room.ID, "roomname", room.Name)

	response = res.NewResponse(res.SuccessResponse{RequestCode: res.CodeType(r.Code())}, logger)
	return
}

func (r KickPeerRequest) Code() CodeType {
	return KICK_PEER
}

func createKickPeerRequest(payload json.RawMessage) (r KickPeerRequest, err error) {
	err = json.Unmarshal(payload, &r)
	return
}
