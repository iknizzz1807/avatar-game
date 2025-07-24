package main

// Import thu vien de dung websocket
import (
	"fmt"
	"net/http"

	"ikniz/avatar/handlers" // Thu vien de xu ly cac message

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	// Thu vien de xu ly JSON
)

func main() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Cho phép tất cả các origin để test
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Serve / roi nang cap len websocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Tao ket noi websocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading connection:", err)
			return
		}

		// Tao mot ket noi moi, voi ID duy nhat
		// Co the su dung uuid.New().String() de tao ID duy nhat
		fmt.Println("WebSocket connection established")

		// Xu ly ket noi
		handlers.HandleConnection(conn, uuid.New().String())
	})

	// Bat dau server
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)

}
