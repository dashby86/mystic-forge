package models

type Enemy struct {
	Name      string
	CurrentHP int
	MaxHP     int
	Attack    int
	Defense   int
	Speed     float64
	Crit      int
	Dodge     int
	Block     int
}
