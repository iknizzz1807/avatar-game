package chat

import "time"

// Thiết kế type cho chat message
type Chat struct {
	UserID   string    `json:"user_id"`
	Message  string    `json:"message"`
	TimeSent time.Time `json:"time_sent"` // Thời gian gửi tin nhắn
}
