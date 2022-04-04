# Transfer Room Ownership
## Description
Change the room owner.

*Note: only the room owner can make this request.*

## Request

```json
{
	"code": "COW",
	"payload": {
		"newOwnerId": "<userID>"
	}
}
```

## Return value

```json
{
	"code": "COW",
	"payload": null
}
```

The updated list of current users in the room is also sent. See [room list of peers](../responses/newpeers.md).


## Examples

```json
// Room owner request to transfer ownership
{
	"code": "COW",
	"payload": {
		"newOwnerId": "c6e39aa0963901a6b347233880b44133647ecd65"
	}
}
// Server response
{
	"code": "COW",
	"payload": null
}
// Updated list of users in room
{
	…
}
```
