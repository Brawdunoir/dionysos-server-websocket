# Create Room
## Description
Create a new room.

## Request

```json
{
	"code": "NRO",
	"payload": {
		"roomname": "<roomname>", // Should be between 3 and 20 characters long
		"isPrivate": <boolean>
	}
}
```

## Return value

```json
{
	"code": "RCS",
	"payload": {
		"roomId": "<roomID>",
		"roomName": "<roomname>"
	}
}
```

The room ID will be mandatory for other users to join this room.

## Examples
### Success

```json
// Client request to create a new room
{
	"code": "NRO",
	"payload": {
		"roomname": "Cinéma",
		"isPrivate": false
	}
}
// Server response
{
	"code": "RCS",
	"payload": {
		"roomId": "cb1f600f162e1730627432b5e093b877b3572780",
		"roomName": "Cinéma"
	}
}
```
### Failure
```json
// Client request to create a new room
{
	"code": "NRO",
	"payload": {
		"roomname": "Ci",
		"isPrivate": false
	}
}
// Server response
{
	"code": "ERR",
	"payload": {
		"error": "RoomName: less than min"
	}
}
```
