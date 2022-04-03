# Join Room
## Description
In case of a *public* room, directly join the room.  
In case of a *private* room, asks the room owner instead and wait for its response.

## Request

```json
{
	"code": "JRO",
	"payload": {
		"roomid": "<roomID>"
	}
}
```

## Return value

```json
{
	"code": "JRO",
	"payload": {
		"roomName": "<roomname>",
		"roomId": "<roomID>",
		"isPrivate": <boolean>
	}
}
```

The list of current users in the room is also sent. See [new peers](../responses/newpeers.md).

## Examples
### Public room

```json
// Client request
{
	"code": "JRO",
	"payload": {
		"roomid": "74b6c5e5af585b04bd606fbd5d458c9072688d4b"
	}
}
// Server response
{
	"code": "JRO",
	"payload": {
		"roomName": "Cinéma",
		"roomId": "74b6c5e5af585b04bd606fbd5d458c9072688d4b",
		"isPrivate": false
	}
}
// List of users in room
{
	…
}
```

### Private room

```json
// Client request
{
	"code": "JRO",
	"payload": {
		"roomid": "e2880f74bf6f46c9db2e7b8011167e1e6cccb3c0"
	}
}
// Server response
{
	"code": "JRO",
	"payload": {
		"roomName": "Cinéma",
		"roomId": "e2880f74bf6f46c9db2e7b8011167e1e6cccb3c0",
		"isPrivate": true
	}
}
// ... Wait room's owner response
```

If the owner *accepts*:
```json
// Server message after room's owner approval
// List of users in room
{
	…
}
```

If the owner *refuses*:
```json
// Server message after room's owner denial
{
	"code": "DEN",
	"payload": {
		"requestCode": "JRO"
	}
}
```


### Failure
```json
// Client request
{
	"code": "JRO",
	"payload": {
		"roomid": "random"
	}
}
// Server response
{
	"code": "ERR",
	"payload": {
		"error": "the room does not exist"
	}
}
```
