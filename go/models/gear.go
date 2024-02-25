package models

import "fmt"

type Gear struct {
	Level   int
	HP      int
	Attack  int
	Defense int
	Speed   int
	Crit    float64
	Dodge   float64
	Block   float64
	SlotId  int
	Rarity  int
}

func (g Gear) GetRarity(rarity int) string {
	fmt.Println("rarity: ", rarity)
	switch rarity {
	case 0:
		return "Junk"
	case 1:
		return "Common"
	case 2:
		return "Uncommon"
	case 3:
		return "Rare"
	case 4:
		return "Epic"
	default:
		return "Not Found"

	}
}
