package enemy

import "math/rand"

type Enemy struct {
	Name    string
	Level   int
	HP      int
	Attack  int
	Defense int
	Speed   int
	Crit    int
	Dodge   int
	Block   int
}

func (e Enemy) CreateEnemy(level int, baseEnemy Enemy) Enemy {
	enemy := baseEnemy
	enemy.Level = level

	// Linear Scaling:
	enemy.HP = baseEnemy.HP + (level-1)*20          // +20 HP per level
	enemy.Attack = baseEnemy.Attack + (level-1)*4   // +4 Attack per level
	enemy.Defense = baseEnemy.Defense + (level-1)*2 // +2 Defense per level
	enemy.Speed = baseEnemy.Defense + (level-1)*1   // +1 Speed per level
	rand.Seed(int64(level))                         // Seed the PRNG

	// Example for Crit (adjust ranges and weights as needed)
	critRange := 0.15 // 15% difference between min and max at level 1
	critWeight := 0.7 // 70% chance of a higher crit value

	crit := rand.Float64() * critRange
	if rand.Float64() < critWeight {
		crit += 0.05 // Bias towards higher end of the range
	}

	// Diminishing Returns after 25%
	if crit > 0.25 {
		diminishingFactor := 0.5 // Adjust this for the severity of the curve
		crit = 0.25 + (crit-0.25)*diminishingFactor
	}

	enemy.Crit = int(crit * 100)

	return enemy
}
