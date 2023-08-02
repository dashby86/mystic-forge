package forge

import (
	"fmt"
	"math/rand"
	"mf/models"
	sqlService "mf/services/sql"
	"time"
)

type Forge struct {
	Sql    sqlService.SqlService
	Player models.Player
}

func (f Forge) CraftGear() models.Gear {
	rand.Seed(time.Now().UnixNano())
	gear := models.Gear{
		Level:   1,
		HP:      determineBaseStats(20),
		Attack:  determineBaseStats(4),
		Defense: determineBaseStats(2),
		Speed:   determineBaseStats(1),
		Crit:    rand.Intn(2) + 1,
		Dodge:   rand.Intn(2) + 1,
		Block:   rand.Intn(2) + 1,
		SlotId:  rand.Intn(8) + 1,
		Rarity:  rand.Intn(2) + 1,
	}
	return gear
}

func determineBaseStats(baseStat int) int {
	rand.Seed(time.Now().UnixNano())
	level := 55
	min := level * 4 * baseStat
	max := level * 6 * baseStat
	return (rand.Intn(max-min+1) + min)
	//return baseStat + (level *)
}

func (f Forge) EquipGear(gear models.Gear) {
	_, err := f.Sql.SaveGearToSlot(f.Player, gear)
	if err != nil {
		return
	}
}

func (f Forge) GenerateRarity(forgeLevel int) string {
	rarityWeights := []float64{
		0.55,  // Junk
		0.30,  // Common
		0.15,  // Uncommon
		0.01,  // Rare
		0.005, // Epic
		//0.0002, // Legendary
		//0.0001, // Mythic
	}

	if forgeLevel >= 10 {
		rarityWeights = append(rarityWeights, 0.0002) // Legendary
	}

	if forgeLevel >= 20 {
		rarityWeights = append(rarityWeights, 0.0001) // Mythic
	}

	rarityIndex := rand.Float64() * 1.0003
	forgeLevel = float64(forgeLevel)

	probability := (forgeLevel-10)*0.001 + 0.005

	for i := len(rarityWeights) - 1; i >= 0; i-- {
		weight := float64(rarityWeights[i])
		if rarityIndex <= (weight + probability) {
			rarity := rarityNames[i]
			prob := weight / 1.0003
			fmt.Println("Crafted - Forge Level:", forgeLevel, "Rarity:", rarity, "Probability:", prob)

			return rarity
		}
	}

	fmt.Println("Crafted - Forge Level:", forgeLevel, "Rarity: Junk", "Probability: Default")
	return "Junk"
}

var rarityNames = []string{
	"Junk",
	"Common",
	"Uncommon",
	"Rare",
	"Epic",
	"Legendary",
	"Mythic",
}
