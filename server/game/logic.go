package game

import (
	"ikniz/avatar/players"
	"ikniz/avatar/types"
	"sync"
)

var (
	playersMutex sync.RWMutex
	Players      = make(map[string]*players.Player)
)

// PlayerState có dạng như sau:
// type PlayerState struct {
// 	ID       string         `json:"id"`       // ID của người chơi
// 	Name     string         `json:"name"`     // Tên của người chơi
// 	Position types.Position `json:"position"` // Vị trí của người chơi
// }

func AddPlayer(id string, name string) {
	playersMutex.Lock()
	defer playersMutex.Unlock()

	// Kiểm tra xem ID đã tồn tại chưa
	if _, ok := Players[id]; ok {
		// Nếu ID đã tồn tại, không thêm người chơi mới
		return
	}

	// Nếu ID chưa tồn tại, tạo người chơi mới và thêm vào danh sách
	newPlayer := &players.Player{
		ID:       id,
		Name:     name,
		Position: types.Position{X: 0, Y: 0}, // Vị trí khởi tạo
	}

	Players[id] = newPlayer
}

func RemovePlayer(playerID string) {
	playersMutex.Lock()
	defer playersMutex.Unlock()
	// Xóa người chơi khỏi map
	delete(Players, playerID)
}

func MovePlayer(playerID string, direction string) {
	playersMutex.Lock()
	defer playersMutex.Unlock()

	player, ok := Players[playerID]
	if !ok {
		return // Người chơi không tồn tại
	}

	// Direction bao gồm "up", "down", "left", "right"
	// Dựa vào đó cập nhật vị trí của người chơi
	switch direction {
	case "up":
		player.UpdatePosition(player.GetPosition().X, player.GetPosition().Y-1)
	case "down":
		player.UpdatePosition(player.GetPosition().X, player.GetPosition().Y+1)
	case "left":
		player.UpdatePosition(player.GetPosition().X-1, player.GetPosition().Y)
	case "right":
		player.UpdatePosition(player.GetPosition().X+1, player.GetPosition().Y)
	}
}
