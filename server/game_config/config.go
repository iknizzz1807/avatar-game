package gameconfig

type Config struct {
	// Thời gian hiển thị bong bóng chat (tính bằng giây)
	ChatBubbleDuration int `json:"chat_bubble_duration"`
}

var GameConfig Config = Config{
	ChatBubbleDuration: 5,
}
