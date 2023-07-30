package models

type Player struct {
	Id        int
	Name      string
	CurrentHP int
	MaxHP     int
	Attack    int
	Defense   int
	Speed     float64
	Crit      int
	Dodge     int
	Block     int
	Ore       int
}
