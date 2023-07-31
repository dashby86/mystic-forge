package forge

import (
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
	return models.Gear{
		Level:   1,
		HP:      determineBaseStats(20),
		Attack:  determineBaseStats(4),
		Defense: determineBaseStats(2),
		Speed:   rand.Float64() + 1,
		Crit:    rand.Float64() + 1,
		Dodge:   rand.Float64() + 1,
		Block:   rand.Float64() + 1,
		SlotId:  rand.Intn(8) + 1,
		Rarity:  rand.Intn(2) + 1,
	}
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
