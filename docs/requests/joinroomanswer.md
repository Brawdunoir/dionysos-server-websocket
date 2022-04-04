# Answer Join Room

## Description
Accept or deny user access to a room.
Should follow [Answer Join Room Pending](../responses/joinroompending.md) response.
Only needed for private rooms.

*Note: only the room owner can make this request.*
## Request

```json
{
	"code": "JRA",
	"payload": {
		"requesterId": "<userID>",
		"accepted": <boolean>
	}
}
```

## Return value

```json
{
	"code": "JRA",
	"payload": null
}
```

The updated list of current users in the room is also sent. See [room list of peers](../responses/newpeers.md).

## Examples

```json
// Server sent answer pending response
{
	"code": "JRP",
	"payload": {
		"roomId": "a056160a7503b6f0df508c2e5b440eb3beba6234",
		"requesterUsername": "Yann",
		"requesterId": "c6e39aa0963901a6b347233880b44133647ecd65"
	}
}
// Client request to accept user Yann
{
	"code": "JRA",
	"payload": {
		"requesterid": "c6e39aa0963901a6b347233880b44133647ecd65",
		"accepted": true
	}
}
// Server response
{
	"code": "JRA",
	"payload": null
}
// Updated list of users in room
{
	…
}
```
