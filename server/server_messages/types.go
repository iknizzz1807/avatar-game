package server_messages

import (
	"ikniz/avatar/types"
	"time"
)

// Init messages for the game logic package
type InitMessage struct {
	Type     string `json:"type"` // "init"
	PlayerID string `json:"player_id"`
}

type ChatMessage struct {
	Type     string    `json:"type"` // "chat"
	UserID   string    `json:"user_id"`
	Message  string    `json:"message"`
	TimeSent time.Time `json:"time_sent"` // Thời gian gửi tin nhắn
}

type MoveMessage struct {
	Type     string         `json:"type"`     // "move"
	Position types.Position `json:"position"` // Position is a struct with X and Y coordinates
}

type GameStateMessage struct {
	Type      string      `json:"type"`       // "game_state"
	GameState interface{} `json:"game_state"` // GameState is a struct that contains the current state of the game
}

// Type của tin nhắn lỗi
type ErrorMessage struct {
	Type    string `json:"type"`    // "error"
	Message string `json:"message"` // Nội dung lỗi
}
