# Transfer File Chunk

## Description
Send file chunks to the room.
This request aims to transfer a piece of file to other peers in the same room.

The writer is responsible to send a chunk of the correct size (5MB). Only the last piece of chunk of the file can be of a lower size.

*Note: only the room owner can make this request.*

## Request

```json
{
	"code": "UCH",
	"payload": {
		"chunkNumber": "<uint>",
	  "chunkData": "<array of bytes>"
  }
}
```

## Return value

```json
{
	"requestCode": "UCH",
	"payload": null
}
```

The file chunk is sent to all room peers, including you. See [this](../responses/uploadchunk.md).

## Examples

```json
// Client request
{
	"code": "UCH",
	"payload": {
		"chunkNumber": "10",
		"chunkData": "[<raw data>]"
	}
}
// Server response
{
	"code": "UCH",
	"payload": null
}
// File chunk returned by server
{
	…
}
```
