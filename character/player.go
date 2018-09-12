package character

import (
	"fmt"
)

type player struct {
	ID        string
	Name      string
	Health    int
	Damage    int
	Intellect int
	Armor     int
}

func (p *player) TakeDamage(amount int) {
	p.Health = p.Health - (amount - p.Armor)
}

func (p player) String() string {
	return fmt.Sprintf("Name: %s\tHealth: %d", p.Name, p.Health)
}
