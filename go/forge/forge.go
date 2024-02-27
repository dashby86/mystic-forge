package forge

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
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
	//calculateLevelAndExpRequired(2000)
	//rand.Seed(time.Now().UnixNano())
	//random level for testing
	playerLevel := f.Player.Level
	level := generateLevel(playerLevel)
	fmt.Println("player level: ", playerLevel, "   level: ", level)
	tier := f.GenerateRarity(f.Player.ForgeLevel)
	crit, dodge, block := calculateSpecializedStat(level, tier)

	gear := models.Gear{
		Level:   level,
		HP:      determineBaseStats(20, level, tier),
		Attack:  determineBaseStats(4, level, tier),
		Defense: determineBaseStats(2, level, tier),
		Speed:   determineBaseStats(1, level, tier),
		Crit:    crit,
		Dodge:   dodge,
		Block:   block,
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
	tierMultipler := tierFloat/5 + 1
	min := int(math.Round(level * 4 * float64(baseStat) * tierMultipler))
	max := int(math.Round(level * 5 * float64(baseStat) * tierMultipler))
	fmt.Printf("min %2f * 4 * %d * (%2f) = %d\n", level, baseStat, tierMultipler, min)
	//fmt.Printf("max %2f * 5 * %d * (%2f) = %d\n", level, baseStat, tierMultipler, max)
	fmt.Printf("tier multiplier %2f\n", tierMultipler)
	//fmt.Printf("tier multiplier %2f\n", tierMultipler)
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
	//forgeLevel = 20
	maxUnlockedTier := (forgeLevel + 15) / 5 // Note: Subtract 1 since unlocking starts at 0

	spew.Dump(rarityWeights)
	fmt.Println("HRERERERE ", maxUnlockedTier+1)
	activeRarityWeights := make([]float64, maxUnlockedTier+1)
	copy(activeRarityWeights, rarityWeights[:maxUnlockedTier+1])
	probabilityIncrease := 0.0025 // 0.025% per forge level
	//forgeLevelOffset := (float64(forgeLevel) - 1) * probabilityIncrease

	// Apply the offset to active weights
	for i := range activeRarityWeights {
		fmt.Printf("%s base: %2f\n", rarityNames[i], activeRarityWeights[i])
		tierUnlockLevel := unlockLevel[i] // Calculate the forge level at which this tier unlocks
		levelOffset := float64(forgeLevel-tierUnlockLevel) * (probabilityIncrease * tierMultiplier[i])
		if levelOffset > 0 { // Apply only if forge level is higher than unlock level
			fmt.Printf("ADDING %2f + %2f \n\n", activeRarityWeights[i], levelOffset)
			activeRarityWeights[i] = activeRarityWeights[i] + levelOffset
		} else {
			fmt.Printf("no its mot higher %d\n\n", i)
		}
	}

	for i := len(activeRarityWeights) - 1; i >= 0; i-- {
		rarityIndex := rand.Float64()
		weight := activeRarityWeights[i]
		fmt.Printf("roll: %2f weight %2f\n", rarityIndex, weight)
		if rarityIndex <= (weight) {
			rarity := rarityNames[i]
			fmt.Println("Crafted - Forge Level:", forgeLevel, "Rarity:", rarity, "Probability:", weight)
			return i
		}
		fmt.Println("Forge Level:", forgeLevel, "Rarity:", rarityNames[i], "Probability:", (weight))
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
	"Ancient",
	"Cosmic",
	"Divine",
	"Primordial",
	"Celestial",
}

var rarityWeights = []float64{
	1.00,   // Junk
	0.45,   // Common
	0.25,   // Uncommon
	0.15,   // Rare
	0.05,   // epic
	0.01,   // legendary
	0.005,  // mythic
	0.005,  // Ancient
	0.001,  // Cosmic
	0.001,  // Divine
	0.001,  // Primordial
	0.0005, // Celestial
}

var tierMultiplier = []float64{
	1,  // Junk
	10, // Common
	7,  // Uncommon
	3,  // Rare
	3,  // epic
	2,  // legendary
	2,  // mythic
	1,  // Ancient
	1,  // Cosmic
	1,  // Divine
	1,  // Primordial
	1,  // Celestial
}

var unlockLevel = []int{
	1,  // Junk
	1,  // Common
	1,  // Uncommon
	1,  // Rare
	5,  // epic
	10, // legendary
	15, // mythic
	20, // Ancient
	25, // Cosmic
	30, // Divine
	35, // Primordial
	40, // Celestial
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
	level := odds[rand.Intn(10)+1]
	if level < 1 {
		level = 1
	}
	return level
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

func calculateSpecializedStat(gearLevel int, tier int) (float64, float64, float64) {
	// ... (Your allocation logic) ...
	primaryStatIndex := rand.Intn(3)
	allocationPool := float64(tier) * 1 * 1.75
	var crit float64
	var dodge float64
	var block float64

	switch primaryStatIndex {
	case 0:
		crit = allocationPool * 0.7
	case 1:
		dodge = allocationPool * 0.7
	case 2:
		block = allocationPool * 0.7
	}

	if rand.Float64() < 0.8 { // 80% chance for a single secondary stat
		secondaryStatIndex := rand.Intn(3)
		multiplier := 0.3
		if secondaryStatIndex == primaryStatIndex {
			multiplier = .5
		}

		switch secondaryStatIndex {
		case 0:
			crit = allocationPool * multiplier
		case 1:
			dodge = allocationPool * multiplier
		case 2:
			block = allocationPool * multiplier
		}
	} else {
		remainingPool := allocationPool * 0.3
		secondaryStatIndex1 := rand.Intn(3)
		for secondaryStatIndex1 == primaryStatIndex {
			secondaryStatIndex1 = rand.Intn(3) // Ensure indices are different
		}
		secondaryStatIndex2 := 3 - primaryStatIndex - secondaryStatIndex1 // Third remaining index

		secondaryStat1Share := remainingPool * rand.Float64() // Random split
		secondaryStat2Share := remainingPool - secondaryStat1Share

		// Assign the values based on the calculated indices
		if secondaryStatIndex1 == 0 {
			crit = secondaryStat1Share
		} else if secondaryStatIndex1 == 1 {
			dodge = secondaryStat1Share
		} else {
			block = secondaryStat1Share
		}

		if secondaryStatIndex2 == 0 {
			crit = secondaryStat2Share
		} else if secondaryStatIndex2 == 1 {
			dodge = secondaryStat2Share
		} else {
			block = secondaryStat2Share
		}
	}

	return crit, dodge, block
}
