// Package này có nghĩa vụ nhận msg từ client -> xử lý, việc xử lý không bao gồm gửi tin nhắn
// mà sẽ gọi package `broadcaster` để gửi tin nhắn đến các client
package handlers

import (
	"encoding/json"
	"fmt"
	"ikniz/avatar/broadcaster"
	"ikniz/avatar/chat"
	client_messages "ikniz/avatar/client_messages"
	"ikniz/avatar/game"
	server_messages "ikniz/avatar/server_messages"
	"ikniz/avatar/types"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Biến quản lý danh sách các kết nối, dùng map để truy cập O(1)
var (
	connections      = make(map[string]*types.Connection)
	connectionsMutex sync.RWMutex
)

// BaseMessage dùng để xác định `type` của message trước khi unmarshal toàn bộ
type BaseMessage struct {
	Type string `json:"type"`
}

func HandleConnection(conn *websocket.Conn, id string) {
	defer conn.Close()

	// Tạo connection instance
	connection := &types.Connection{Ws: conn, PlayerID: id}

	// Thêm vào danh sách connections (thread-safe)
	connectionsMutex.Lock()
	connections[id] = connection
	connectionsMutex.Unlock()

	// Cleanup khi connection đóng
	defer func() {
		// Remove connection khỏi danh sách (thread-safe)
		connectionsMutex.Lock()
		delete(connections, id)
		connectionsMutex.Unlock()

		// Remove người chơi khỏi game
		game.RemovePlayer(id)

		// Broadcast lại game state sau khi người chơi thoát
		broadcaster.BroadcastGameState(getAllConnections())
		fmt.Printf("Connection %s removed\n", id)
	}()

	for {
		// Đọc message dưới dạng byte array để xử lý hiệu quả hơn
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			// Chỉ log lỗi nếu đó không phải là lỗi đóng kết nối bình thường
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Error reading message from client %s: %v\n", id, err)
			}
			break
		}

		// Unmarshal vào struct cơ bản để lấy type
		var baseMsg BaseMessage
		if err := json.Unmarshal(messageBytes, &baseMsg); err != nil {
			fmt.Println("Error unmarshalling base message:", err)
			continue
		}

		// Switch case theo type để unmarshal vào struct cụ thể
		switch baseMsg.Type {
		case "init":
			var initMsg client_messages.InitMessage
			if err := json.Unmarshal(messageBytes, &initMsg); err != nil {
				fmt.Println("Error unmarshalling init message:", err)
				continue
			}
			HandleInitMessage(connection, initMsg)

		case "move":
			var moveMsg client_messages.MoveMessage
			if err := json.Unmarshal(messageBytes, &moveMsg); err != nil {
				fmt.Println("Error unmarshalling move message:", err)
				continue
			}
			HandleMoveMessage(connection, moveMsg)

		case "chat":
			var chatMsg client_messages.ChatMessage
			if err := json.Unmarshal(messageBytes, &chatMsg); err != nil {
				fmt.Println("Error unmarshalling chat message:", err)
				continue
			}
			HandleChatMessage(connection, chatMsg)

		default:
			fmt.Printf("Unknown message type: %s\n", baseMsg.Type)
		}
	}
}

// Hàm xử lý message init
// Nhận vào biến `conn` là một con trỏ đến `Connection` để có thể truy cập các trường như `ws` và `PlayerID`
// Nhận vào biến `msg` là một struct `InitMessage` đã được unmarshal từ JSON, đây là message từ client trong lúc đăng nhập vào server
func HandleInitMessage(conn *types.Connection, msg client_messages.InitMessage) {
	fmt.Printf("Handling init message: %+v\n", msg)

	// Gui tin nhan initialization den client
	initMessageServer := server_messages.InitMessage{
		Type:     "init",
		PlayerID: conn.PlayerID,
	}

	// Gửi message khởi tạo từ server đến client
	err := broadcaster.SendMessageToClient(conn, websocket.TextMessage, initMessageServer)
	if err != nil {
		fmt.Println("Error sending init message to client:", err)
		return
	}

	playerID := conn.PlayerID
	playerName := msg.PlayerName

	if playerID == "" || playerName == "" {
		message := "Player ID or Player Name cannot be empty"
		fmt.Println("Error:", message)
		// Gửi thông báo lỗi đến client
		broadcaster.SendErrorMessageToClient(conn, message)
		return
	}

	// Them nguoi choi vao game
	game.AddPlayer(playerID, playerName)
	broadcaster.BroadcastGameState(getAllConnections())
}

// Hàm xử lý message move
func HandleMoveMessage(conn *types.Connection, msg client_messages.MoveMessage) {
	fmt.Printf("Handling move message: %+v\n", msg)
	// Xử lý logic di chuyển, cập nhật vị trí của người chơi
	game.MovePlayer(conn.PlayerID, msg.Direction)
	broadcaster.BroadcastGameState(getAllConnections())
}

// Hàm xử lý message chat
func HandleChatMessage(conn *types.Connection, msg client_messages.ChatMessage) {
	fmt.Printf("Handling chat message: %+v\n", msg)

	// Xử lý logic chat
	chatMessage := chat.Chat{
		UserID:   conn.PlayerID,
		Message:  msg.Message,
		TimeSent: time.Now(),
	}

	chat.AddChat(chatMessage)
	broadcaster.BroadcastGameState(getAllConnections())
}

// Hàm helper để lấy danh sách tất cả các connection một cách an toàn
func getAllConnections() []*types.Connection {
	connectionsMutex.RLock()
	defer connectionsMutex.RUnlock()

	conns := make([]*types.Connection, 0, len(connections))
	for _, c := range connections {
		conns = append(conns, c)
	}
	return conns
}
