# Room List of Peers
## Description
A change in the peers has occured in your room.

*Note: This can be anything from a peer rename to a kick or a new peer joining your room.*

## Response

```json
{
	"code": "NEP",
	"payload": {
		"ownerId": "<userID>",
		"peers": [
      {
        "id": "<userID>",
        "name": "<username>"
      },

      …

      {
        "id": "<userID>",
        "name": "<username>"
      },
    ]
	}
}
```

## Examples
```json
{
	"code": "NEP",
	"payload": {
		"ownerId": "41805c7077a15e592810a0495072844b2cd72c8c",
		"peers": [{
			"id": "41805c7077a15e592810a0495072844b2cd72c8c",
			"name": "Yann"
		}, {
			"id": "c6e39aa0963901a6b347233880b44133647ecd65",
			"name": "Romain"
		}]
	}
}
```
