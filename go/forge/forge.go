package forge

import (
	"fmt"
	"log"
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
	calculateLevelAndExpRequired(2000)
	rand.Seed(time.Now().UnixNano())
	//random level for testing
	playerLevel := rand.Intn(50) + 1
	level := generateLevel(playerLevel)
	fmt.Println("player level: ", playerLevel, "   level: ", level)
	tier := f.GenerateRarity(f.Player.ForgeLevel)
	gear := models.Gear{
		Level:   level,
		HP:      determineBaseStats(20, level, tier),
		Attack:  determineBaseStats(4, level, tier),
		Defense: determineBaseStats(2, level, tier),
		Speed:   determineBaseStats(1, level, tier),
		Crit:    rand.Intn(2) + 1,
		Dodge:   rand.Intn(2) + 1,
		Block:   rand.Intn(2) + 1,
		SlotId:  rand.Intn(8) + 1,
		Rarity:  tier,
	}
	err := f.Sql.GrantExp(f.Player.Id, tier+1)
	if err != nil {
		log.Fatal(err)
	}
	return gear
}

func determineBaseStats(baseStat int, gearLevel int, tier int) int {
	rand.Seed(time.Now().UnixNano())
	level := float64(gearLevel)
	tierFloat := float64(tier)
	tierMultipler := tierFloat/10 + 1
	min := int(math.Round(level * 4 * float64(baseStat) * tierMultipler))
	max := int(math.Round(level * 5 * float64(baseStat) * tierMultipler))
	fmt.Printf("min %2f * 4 * %d * (%2f) = %d\n", level, baseStat, tierMultipler, min)
	fmt.Printf("max %2f * 5 * %d * (%2f) = %d\n", level, baseStat, tierMultipler, max)
	fmt.Printf("tier multiplier %2f\n", tierMultipler)
	fmt.Printf("tier multiplier %2f\n", tierMultipler)
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
	forgeLevel = 20
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

	forgeLevelFloat := float64(forgeLevel)
	probability := (forgeLevelFloat-10)*0.001 + 0.005
	for i := len(rarityWeights) - 1; i >= 0; i-- {
		rarityIndex := rand.Float64() * 1.0003
		weight := float64(rarityWeights[i])
		fmt.Printf("roll: %2f weight %2f\n", rarityIndex, weight+probability)
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

func generateLevel(playerLevel int) int {
	odds := []int{
		playerLevel - 2,
		playerLevel - 1,
		playerLevel - 1,
		playerLevel,
		playerLevel,
		playerLevel,
		playerLevel,
		playerLevel + 1,
		playerLevel + 1,
		playerLevel + 2,
		playerLevel + 3,
	}
	return odds[rand.Intn(10)+1]
}

func calculateLevelAndExpRequired(accumulatedExp float64) (int, float64) {
	baseExp := 20.00
	expIncrement := 0.10

	level := 1.0
	expRequired := baseExp

	for accumulatedExp >= expRequired {
		level++
		expRequired = expRequired * (1 + expIncrement)
		fmt.Println("weird calc: ", math.Round(expRequired*(1.00+expIncrement)))
	}

	fmt.Printf("level: %2f\n", level)

	return int(level), expRequired - accumulatedExp
}

func (f Forge) grantExp(exp int) {

}
