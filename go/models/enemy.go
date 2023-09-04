package models

type Enemy struct {
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
}

// Implementing the Character interface for Enemy
func (e *Enemy) GetId() int {
	return e.Id
}

func (e *Enemy) GetName() string {
	return e.Name
}

func (e *Enemy) GetCurrentHP() int {
	return e.CurrentHP
}

func (e *Enemy) GetMaxHP() int {
	return e.MaxHP
}

func (e *Enemy) GetAttack() int {
	return e.Attack
}

func (e *Enemy) GetDefense() int {
	return e.Defense
}

func (e *Enemy) GetSpeed() float64 {
	return e.Speed
}

func (e *Enemy) GetCrit() float64 {
	return e.Crit
}

func (e *Enemy) GetDodge() float64 {
	return e.Dodge
}

func (e *Enemy) GetBlock() float64 {
	return e.Block
}

func (e *Enemy) Defend(damage int) {
	// Apply the damage to the *Enemy.
	e.CurrentHP -= damage
}
