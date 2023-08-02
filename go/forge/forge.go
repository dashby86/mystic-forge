package forge

import (
	"fmt"
	"math"
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
	tier := f.GenerateRarity(f.Player.ForgeLevel)
	playerLevel := rand.Intn(50) + 1
	gear := models.Gear{
		Level:   playerLevel,
		HP:      determineBaseStats(20, playerLevel, tier),
		Attack:  determineBaseStats(4, playerLevel, tier),
		Defense: determineBaseStats(2, playerLevel, tier),
		Speed:   determineBaseStats(1, playerLevel, tier),
		Crit:    rand.Intn(2) + 1,
		Dodge:   rand.Intn(2) + 1,
		Block:   rand.Intn(2) + 1,
		SlotId:  rand.Intn(8) + 1,
		Rarity:  tier,
	}
	return gear
}

func determineBaseStats(baseStat int, playerLevel int, tier int) int {
	rand.Seed(time.Now().UnixNano())
	level := float64(playerLevel)
	tierFloat := float64(tier)
	tierMultipler := tierFloat/10 + 1
	min := int(math.Round(level * 4 * float64(baseStat) * tierMultipler))
	max := int(math.Round(level * 5 * float64(baseStat) * tierMultipler))
	fmt.Printf("min %d * 4 * %d * (%d) = %d\n", level, baseStat, tierMultipler, min)
	fmt.Printf("max %d * 5 * %d * (%d) = %d\n", level, baseStat, tierMultipler, max)
	fmt.Printf("tier multiplier %d\n", tierMultipler)
	fmt.Printf("tier multiplier %d\n", tierMultipler)
	return rand.Intn(max-min+1) + min
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
func (f Forge) GenerateRarity(forgeLevel int) int {
	rarityWeights := []float64{
		0.55,  // Junk
		0.30,  // Common
		0.15,  // Uncommon
		0.05,  // Rare
		0.005, // Epic
	}
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
			return i
		}
		fmt.Println("Forge Level:", forgeLevel, "Rarity:", rarityNames[i], "Probability:", (weight+probability)/1.0003)
	}
	fmt.Println("Crafted - Forge Level:", forgeLevel, "Rarity: Junk", "Probability: Default")
	return 0
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
