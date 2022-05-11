# Transfer File Metadata

## Description
File metadata to be able to accept further incoming chunks of file.

## Response
	Name      string `json:"name"`      // Name of file
	Size      uint64 `json:"size"`      // Size of file, in bytes
	ChunkSize uint   `json:"chunkSize"` // Chunk size, in bytes
	Chunks    uint16 `json:"chunks"`    // Number of chunks

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
