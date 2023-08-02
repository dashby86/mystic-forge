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
	tier := f.GenerateRarity(f.Player.)
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

func determineBaseStats(baseStat int, playerLevel int, tier int) int {
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

/**
* todo pull tiers from db
*
 */
func (f Forge) GenerateRarity(forgeLevel int) string {
	rarityWeights := []float64{
		0.55,  // Junk
		0.30,  // Common
		0.15,  // Uncommon
		0.05,  // Rare
		0.005, // Epic
	}
	forgeLevel = rand.Intn(50) + 1
	if forgeLevel >= 10 {
		rarityWeights = append(rarityWeights, 0.0002) // Legendary
	}
	if forgeLevel >= 20 {
		rarityWeights = append(rarityWeights, 0.0001) // Mythic
	}
	rand.Seed(time.Now().UnixNano())
	rarityIndex := rand.Float64() * 1.0003

	forgeLevelFloat := float64(forgeLevel)
	probability := (forgeLevelFloat-10)*0.001 + 0.005
	for i := len(rarityWeights) - 1; i >= 0; i-- {
		weight := float64(rarityWeights[i])
		fmt.Printf("roll: %d weight %d\n", rarityIndex, weight+probability)
		if rarityIndex <= (weight + probability) {
			rarity := rarityNames[i]
			prob := (weight + probability) / 1.0003
			fmt.Println("Crafted - Forge Level:", forgeLevel, "Rarity:", rarity, "Probability:", prob)
			return rarityNames[i]
		}
		fmt.Println("Forge Level:", forgeLevel, "Rarity:", rarityNames[i], "Probability:", (weight+probability)/1.0003)
	}
	fmt.Println("Crafted - Forge Level:", forgeLevel, "Rarity: Junk", "Probability: Default")
	return rarityNames[0]
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
