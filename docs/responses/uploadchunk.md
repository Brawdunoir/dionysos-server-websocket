# Transfer File Chunk

## Description
Receiv file chunks to the room.
This request aims to transfer a piece of file to other peers in the same room.

The readers (i.e. other peers) is responsible to write the file accordingly to the file metadata transfered before (see [Transfer File Metadata](../responses/loadfile.md)) by moving the writing head.

## Response

The response is the same as the request by the writer (room owner).

```json
{
	"code": "UCH",
	"payload": {
		"chunkNumber": "<uint>",
	  "chunkData": "<array of bytes>"
  }
}
```

## Examples
See full example in [Transfer File Chunk](../requests/uploadFile.md).
