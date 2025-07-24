package types

import "github.com/gorilla/websocket"

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Connection struct {
	Ws       *websocket.Conn
	PlayerID string // Phân biệt các kết nối
}
