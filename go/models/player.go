package models

import "math"

type Player struct {
	Id           int
	Name         string
	Experience   int
	Level        int
	ForgeLevel   int
	DungeonLevel int
	HP           int
	Attack       int
	Defense      int
	Speed        int
	Crit         float64
	Dodge        float64
	Block        float64
	Ore          int
}

func (player Player) CalculateRequiredExp(targetLevel int) int {
	exponent := 1.1
	baseExp := 20

	requiredExp := float64(baseExp) * math.Pow(float64(targetLevel), exponent)
	return int(requiredExp)
}
