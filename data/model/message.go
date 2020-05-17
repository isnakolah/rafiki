package model

import (
	"github.com/google/uuid"
)

type Message struct {
	MessageId       uuid.UUID `gorm:"type:uuid; primary_key" json:"message_id"`
	UserId          uuid.UUID `gorm:"type:uuid; not null" json:"user_id"`
	RecipientNumber string    `gorm:"type:varchar(15);"  json:"recipient_number"`
	MessageBody     string    `gorm:"type:varchar(5000);"  json:"message_body"`
	MessageSent     bool      `gorm:"default:true" json:"message_sent"`
	Response        string    `gorm:"type:varchar(500);"  json:"response"`
	Rafiki
}
