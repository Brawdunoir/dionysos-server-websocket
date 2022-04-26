package objects

// File metadata gathers metadata about a file
type FileMetadata struct {
	Name      string `json:"name"`      // Name of file
	Size      uint64 `json:"size"`      // Size of file, in bytes
	ChunkSize uint   `json:"chunkSize"` // Chunk size, in bytes
	Chunks    uint16 `json:"chunks"`    // Number of chunks
}
