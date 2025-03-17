package types

import "github.com/google/uuid"


type LiveMessage struct {
	MessageType string    `json:"message_type"`
	SentByID    uuid.UUID `json:"sent_by_id"`
	Data        any       `json:"data"`
}