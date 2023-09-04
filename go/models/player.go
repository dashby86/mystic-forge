package models

type Player struct {
	Id        int     `json:"Id"`
	Name      string  `json:"Name"`
	CurrentHP int     `json:"CurrentHP"`
	MaxHP     int     `json:"MaxHP"`
	Attack    int     `json:"Attack"`
	Defense   int     `json:"Defense"`
	Speed     float64 `json:"Speed"`
	Crit      float64 `json:"Crit"`
	Dodge     float64 `json:"Dodge"`
	Block     float64 `json:"Block"`
	Ore       int     `json:"Ore"`
}

// Implementing the Character interface for Player
func (p *Player) GetId() int {
	return p.Id
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) GetCurrentHP() int {
	return p.CurrentHP
}

func (p *Player) GetMaxHP() int {
	return p.MaxHP
}

func (p *Player) GetAttack() int {
	return p.Attack
}

func (p *Player) GetDefense() int {
	return p.Defense
}

func (p *Player) GetSpeed() float64 {
	return p.Speed
}

func (p *Player) GetCrit() float64 {
	return p.Crit
}

func (p *Player) GetDodge() float64 {
	return p.Dodge
}

func (p *Player) GetBlock() float64 {
	return p.Block
}

func (p *Player) GetOre() int {
	return p.Ore
}

func (p *Player) Defend(damage int) {
	// Apply the damage to the enemy.
	p.CurrentHP -= damage
}
