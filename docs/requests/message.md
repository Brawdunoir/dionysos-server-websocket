# Send Message

## Description
Send message in the room chat.

*Note: The messages are just forwarded to the list of peers, a new user joining the room will have an empty chat.*

## Request

```json
{
	"code": "NME",
	"payload": {
		"content": "<message>",
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

The message is also sent to all peers, including yourself.

## Examples

```json
// Client request to send a message
{
	"code": "NME",
	"payload": {
		"content": "Hello World!",
	}
}
// Server acknowledge request
{
	"code": "NME",
	"payload": null
}
// This server response is also sent to all peers, including yourself.
{
	"code": "NME",
	"payload": {
		"senderId": "c6e39aa0963901a6b347233880b44133647ecd65",
		"senderName": "Yann",
		"content": "Hello World!",
	}
}
```
