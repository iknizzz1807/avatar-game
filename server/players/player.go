package players

import (
	"ikniz/avatar/types"
)

func NewPlayer(id string, name string, posX float64, posY float64) *Player {
	return &Player{
		ID:   id,
		Name: name,
		Position: types.Position{
			X: posX,
			Y: posY,
		},
	}
}

func (p *Player) UpdatePosition(posX float64, posY float64) {
	p.Position.X = posX
	p.Position.Y = posY
}

// Get Position returns the player's position as a game.Position struct
func (p *Player) GetPosition() types.Position {
	return p.Position
}
