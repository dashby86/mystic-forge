package models

type Player struct {
	Id         int
	Name       string
	Experience int
	Level      int
	ForgeLevel int
	HP         int
	Attack     int
	Defense    int
	Speed      int
	Crit       float64
	Dodge      float64
	Block      float64
	Ore        int
}
