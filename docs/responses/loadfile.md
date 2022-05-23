# Transfer File Metadata

## Description
File metadata to be able to accept further incoming chunks of file.

## Response

```json
{
	"code": "LFI",
	"payload": {
		"name": "<fileName>",
		"size": "<fileSize>",
		"chunkSize": "<chunkSize>",
		"chunks": "<chunksNumber>"
	}
}
```

## Examples
See full example in [file transfer](../requests/loadfile.md).
