# Transfer File Metadata

## Description
Send file metadata to the room.
This request aims to advertise room peers to know file size and other information related for the further transfer.

*Note: only the room owner can make this request.*

## Request

```json
{
	"code": "LFI",
	"payload": {
		"name": "<fileName", // should be between 3 and 254 characters
		"size": <uint>
	}
}
```

## Return value

```json
{
	"requestCode": "LFI",
	"payload": null
}
```

The file metadata and other informations is sent to all room peers, including you. See [this](../responses/loadfile.md).

## Examples

```json
// Client request
{
	"code": "LFI",
	"payload": {
		"name": "Inception",
		"size": 2641404887 //2.46 GB
	}
}
// Server response
{
	"code": "LFI",
	"payload": null
}
// File metadata processed by server
{
	…
}
```
