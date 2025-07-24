// Package này có nghĩa vụ liên tục gửi trạng thái của game đến client
package broadcaster

import (
	"encoding/json"
	"fmt"
	"ikniz/avatar/game"
	"ikniz/avatar/server_messages"
	"ikniz/avatar/types"

	"github.com/gorilla/websocket"
)

// Hàm để gửi một tin nhắn cụ thể đến một client cụ thể
// messageType là kiểu của tin nhắn mà bạn muốn gửi đến client. Trong trường hợp này, nó có thể là một hằng số từ gói websocket như:
// - websocket.TextMessage: để gửi tin nhắn dạng văn bản (text) = 1
// - websocket.BinaryMessage: để gửi tin nhắn dạng nhị phân (binary) = 2
// - websocket.CloseMessage: để gửi tin nhắn đóng kết nối (close) = 8
func SendMessageToClient(conn *types.Connection, messageType int, message interface{}) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}

	if err := conn.Ws.WriteMessage(messageType, messageBytes); err != nil {
		return fmt.Errorf("error sending message to client: %w", err)
	}
	return nil
}

// Hàm để gửi một tin nhắn cụ thể đến tất cả các client
func SendMessageToAllClients(conn []*types.Connection, messageType int, message interface{}) {
	for _, conn := range conn {
		byteMessage, ok := message.([]byte)
		if !ok {
			fmt.Println("Error: message is not a []byte")
			continue
		}
		if err := conn.Ws.WriteMessage(messageType, byteMessage); err != nil {
			fmt.Println("Error sending message to client:", err)
			continue
		}
	}
}

// BroadcastGameState sends the current game state to all connected clients
func BroadcastGameState(clients []*types.Connection) {
	gameState := game.GetGameState()
	gameStateMessage := server_messages.GameStateMessage{
		Type:      "game_state",
		GameState: gameState,
	}
	// Convert game state to JSON
	gameStateMessageJSON, err := json.Marshal(gameStateMessage)
	if err != nil {
		fmt.Println("Error marshalling game state:", err)
		return
	}

	// Send the game state to all connected clients
	for _, conn := range clients {
		if err := conn.Ws.WriteMessage(websocket.TextMessage, gameStateMessageJSON); err != nil {
			fmt.Println("Error sending game state to client:", err)
			continue
		}
	}
}

// Gửi tin nhắn lỗi đến client, sử dụng kiểu ErrorMessage
func SendErrorMessageToClient(conn *types.Connection, errorMessage string) error {
	errorMsg := server_messages.ErrorMessage{
		Type:    "error",
		Message: errorMessage,
	}

	// Chuyển đổi message thành JSON
	messageBytes, err := json.Marshal(errorMsg)
	if err != nil {
		fmt.Println("Error marshalling error message:", err)
		return fmt.Errorf("error marshalling error message: %w", err)
	}

	// Gửi tin nhắn lỗi đến client
	if err := conn.Ws.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		fmt.Println("Error sending error message to client:", err)
		return fmt.Errorf("error sending error message to client: %w", err)
	}

	return nil
}
