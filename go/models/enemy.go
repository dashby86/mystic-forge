package models

type Enemy struct {
	Name      string
	CurrentHP int
	MaxHP     int
	Attack    int
	Defense   int
	Speed     float64
	Crit      float64
	Dodge     float64
	Block     float64
}
