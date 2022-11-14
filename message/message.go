// Message struct to be passed between client and server
package message

import (
	"github.com/google/uuid"
)

type Message struct {
	Id      string
	Headers map[string]string
	Status  string
	Body    string
}

func NewMessage() Message {
	message := new(Message)
	id := uuid.New()
	message.Id = id.String()
	return *message
}
