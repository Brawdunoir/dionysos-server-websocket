# Answer Join Room Pending

## Description
A request to join your room needs your attention.

*Note: it is a response from the server because of an another user request.*

## Response

```json
{
	"code": "JRP",
	"payload": {
		"roomId": "<roomID>",
		"requesterUsername": "<username>",
		"requesterId": "<userID>"
	}
}
```

## Examples
See full example in [answer join room](../requests/joinroomanswer.md).
