# Kick Peer

## Description
Kick a user from your room.

*Note: only the room owner can make this request.*
## Request

```json
{
	"code": "KPE",
	"payload": {
		"peerId": "<userID>",
	}
}
```

## Return value

```json
{
	"code": "KPE",
	"payload": null
}
```

The updated list of current users in the room is also sent. See [new peers](../responses/newpeers.md).

## Examples

```json
// Client request to kick a user
{
	"code": "KPE",
	"payload": {
		"peerId": "c6e39aa0963901a6b347233880b44133647ecd65",
	}
}
// Server response
{
	"code": "KPE",
	"payload": null
}
// Updated list of users in room
{
	…
}
```
