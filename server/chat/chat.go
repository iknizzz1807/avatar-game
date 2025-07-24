// Quản lý tất cả những gì thuộc về chat trong game
package chat

import "sync"

const maxChatHistory = 100 // Giới hạn số lượng tin nhắn được lưu trữ

var (
	chatsMutex sync.RWMutex
	Chats      []Chat
)

func AddChat(chat Chat) {
	chatsMutex.Lock()
	defer chatsMutex.Unlock()

	Chats = append(Chats, chat)

	// Nếu số lượng tin nhắn vượt quá giới hạn, cắt bớt tin nhắn cũ nhất
	if len(Chats) > maxChatHistory {
		Chats = Chats[len(Chats)-maxChatHistory:]
	}
}

func GetChats() []Chat {
	chatsMutex.RLock()
	defer chatsMutex.RUnlock()

	// Trả về một bản sao của slice để tránh race condition ở phía người gọi
	chatsCopy := make([]Chat, len(Chats))
	copy(chatsCopy, Chats)
	return chatsCopy
}
