# Change Username
## Description
Change the user's username.

## Request

```json
{
	"code": "CHU",
	"payload": {
		"newUsername": "<username>" // Should be between 3 and 20 characters long
	}
}
```

## Return value

```json
{
	"code": "CHU",
	"payload": {
		"username": "<username>"
	}
}
```

## Examples
### Success

```json
// Client request
{
	"code": "CHU",
	"payload": {
		"newUsername": "Yann"
	}
}
// Server response
{
	"code": "CHU",
	"payload": {
		"username": "Yann"
	}
}
```
### Failure
```json
// Client request
{
	"code": "CHU",
	"payload": {
		"newUsername": "Jo"
	}
}
// Server response
{
	"code": "ERR",
	"payload": {
		"error": "NewUsername: less than min"
	}
}
```
