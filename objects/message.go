package objects

// Message is a content associated to a sender by its ID
type Message struct {
	SenderID   string `json:"senderId"`
	SenderName string `json:"senderName"`
	Content    string `json:"content"`
}

// NewMessage creates a new message
func NewMessage(senderID, senderName, content string) Message {
	return Message{SenderID: senderID, SenderName: senderName, Content: content}
}
