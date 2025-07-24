package players

import "ikniz/avatar/types"

type Player struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Position types.Position `json:"position"`
}
