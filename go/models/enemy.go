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
