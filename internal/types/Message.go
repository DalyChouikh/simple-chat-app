package types

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

func NewMessage(sender, recipient, content string) *Message {
	return &Message{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
	}
}
