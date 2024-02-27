package battle

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"math/rand"
	"mf/enemy"
	"mf/models"
	sqlService "mf/services/sql"
	"time"
)

type Battle struct {
	Player models.Player
	Enemy  enemy.Enemy
	Sql    sqlService.SqlService
}

func (battle Battle) SimBattle() {
	// Initialize the turn variable.
	//turn := 1
	enemyHp := battle.Enemy.HP
	playerHp := battle.Player.HP

	spew.Dump(battle.Enemy)

	fmt.Printf("Starting battle! Player life: %d  -- Enemy life: %d\n", playerHp, enemyHp)

	// Simulate the battle.
	for {
		// Check which entity gets to attack first.
		if battle.Player.Speed > battle.Enemy.Speed {
			playerAttacks(battle.Player, battle.Enemy, &playerHp, &enemyHp)
			if enemyHp > 0 {
				enemyAttacks(battle.Player, battle.Enemy, &playerHp, &enemyHp)
			}
		} else {
			enemyAttacks(battle.Player, battle.Enemy, &playerHp, &enemyHp)
			if playerHp > 0 {
				playerAttacks(battle.Player, battle.Enemy, &playerHp, &enemyHp)
			}
		}

		// Check if the battle is over.
		if playerHp <= 0 {
			fmt.Println("The player has died.")
			break
		}
		if enemyHp <= 0 {
			fmt.Println("The enemy has died.")

			currentDungeonLevel := battle.Player.DungeonLevel + 1
			err := battle.Sql.UpdateDungeonLevel(battle.Player.Id, currentDungeonLevel)
			if err != nil {
				// Handle the error
			}

			_, err = battle.Sql.GrantOre(battle.Player.Id, 10)
			if err != nil {
				// Handle the error
			}
			err = battle.Sql.UpdateForgeLevel(battle.Player)
			if err != nil {
				// Handle the error
			}

			break
		}

		fmt.Printf("Player life: %d  -- Enemy life: %d\n", playerHp, enemyHp)

		time.Sleep(2 * time.Second)
	}
}

func playerAttacks(player models.Player, enemy enemy.Enemy, playerHp *int, enemyHp *int) {
	// Calculate the player's damage.
	damage := player.Attack
	baseReduction := float64(enemy.Defense) / (float64(enemy.Defense) + 200.0) // Add .0 to the constant
	// Constant of 100 for now
	mitigatedDamage := int(float64(damage) * baseReduction)
	finalDamage := max(damage-mitigatedDamage, 0)

	if float64(rand.Intn(100)) < player.Crit { // Convert to float64
		fmt.Println("Critical Strike!")
		finalDamage *= 2
	}

	if float64(rand.Intn(100)) < enemy.Dodge { // Convert to float64
		fmt.Println("Enemy Dodged!")
	} else {
		// Apply the damage to the enemy.
		fmt.Printf("Attacking %s for %d damage.\n", enemy.Name, finalDamage)
		*enemyHp -= finalDamage
	}

	fmt.Println("Enemy Defense:", enemy.Defense)
	fmt.Println("Base Reduction:", baseReduction)
	fmt.Println("Mitigated Damage:", mitigatedDamage)

	fmt.Println("\n\n\n")
}

func enemyAttacks(player models.Player, enemy enemy.Enemy, playerHp *int, enemyHp *int) {
	// Calculate the enemy's damage.

	damage := enemy.Attack
	baseReduction := float64(player.Defense) / (float64(player.Defense) + 200.0) // Add .0 to the constant
	// Constant of 100 for now
	mitigatedDamage := int(float64(damage) * baseReduction)
	finalDamage := max(damage-mitigatedDamage, 0)

	if float64(rand.Intn(100)) < enemy.Crit {
		finalDamage *= 2
	}

	if float64(rand.Intn(100)) < player.Dodge {
		fmt.Println("Player Dodged!")
	} else {
		// Apply the damage to the player.
		fmt.Printf("Attacking %s for %d damage.\n", player.Name, finalDamage)
		*playerHp -= finalDamage
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
