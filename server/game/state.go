package game

import (
	"ikniz/avatar/chat"
	"ikniz/avatar/players"
)

// Hàm này được gọi liên tục để gửi trạng thái game đến client
func GetGameState() GameState {
	// Khóa và lấy danh sách người chơi
	playersMutex.RLock()
	playersList := make([]players.Player, 0, len(Players))
	for _, player := range Players {
		playersList = append(playersList, *player)
	}
	playersMutex.RUnlock()

	// Lấy tin nhắn chat (đã được bảo vệ bởi mutex trong package chat)
	chats := chat.GetChats()

	return GameState{
		Players:      playersList,
		ChatMessages: chats,
	}
}
