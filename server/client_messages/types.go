// Định nghĩa type cho các message từ client gửi lên server
// Mục đích để giải mã JSON từ client
package clientmessages

type InitMessage struct {
	Type       string `json:"type"` // "init"
	PlayerName string `json:"player_name"`
}

type MoveMessage struct {
	// Khong can ID vi client da co ID tu server
	// them vao do, ta khong nen tin tuong vao client
	Type      string `json:"type"`      // "move"
	Direction string `json:"direction"` // "up", "down", "left", "right"
}

type ChatMessage struct {
	Type    string `json:"type"` // "chat"
	Message string `json:"message"`
}
