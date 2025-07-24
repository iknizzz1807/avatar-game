package game

import (
	"ikniz/avatar/chat"
	"ikniz/avatar/players"
)

// Struct này sẽ luôn được mở rộng trong tương lai
type GameState struct {
	Players      []players.Player `json:"players"`
	ChatMessages []chat.Chat      `json:"chat"`
}
