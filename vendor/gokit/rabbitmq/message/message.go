package message

type Message struct {
	RouteKey   string
	Content    string
	Expiration string
}

func NewMessage(key, content string) *Message {
	return NewMessageWithExpiration(key, content, "")
}

func NewMessageWithExpiration(key, content, expiration string) *Message {
	return &Message{key, content, expiration}
}
