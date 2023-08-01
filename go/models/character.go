package models

type Character interface {
	GetId() int
	GetName() string
	GetCurrentHP() int
	GetMaxHP() int
	GetAttack() int
	GetDefense() int
	GetSpeed() float64
	GetCrit() float64
	GetDodge() float64
	GetBlock() float64
}
