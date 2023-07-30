package battle

import (
	"fmt"
	"math/rand"
	"mf/models"
	sqlService "mf/services/sql"
	"time"
)

type Battle struct {
	Player models.Player
	Enemey models.Enemy
	Sql    sqlService.SqlService
}

func (battle Battle) SimBattle() {
	// Initialize the turn variable.
	//turn := 1
	enemyHp := battle.Enemey.HP
	playerHp := battle.Player.HP

	fmt.Printf("Starting battle! Player life: %d  -- Enemy life: %d\n", playerHp, enemyHp)

	// Simulate the battle.
	for {
		// Check which entity gets to attack first.
		if battle.Player.Speed > battle.Enemey.Speed {
			playerAttacks(battle.Player, battle.Enemey, &playerHp, &enemyHp)
			if enemyHp > 0 {
				enemyAttacks(battle.Player, battle.Enemey, &playerHp, &enemyHp)
			}
		} else {
			enemyAttacks(battle.Player, battle.Enemey, &playerHp, &enemyHp)
			if playerHp > 0 {
				playerAttacks(battle.Player, battle.Enemey, &playerHp, &enemyHp)
			}
		}

		// Check if the battle is over.
		if playerHp <= 0 {
			fmt.Println("The player has died.")
			break
		}
		if enemyHp <= 0 {
			fmt.Println("The enemy has died.")

			break
		}

		fmt.Printf("Player life: %d  -- Enemy life: %d\n", playerHp, enemyHp)

		time.Sleep(2 * time.Second)
	}
}

func playerAttacks(player models.Player, enemy models.Enemy, playerHp *int, enemyHp *int) {
	// Calculate the player's damage.
	damage := player.Attack - enemy.Defense
	if rand.Intn(100) < player.Crit {
		fmt.Println("Critical Strike!")
		damage *= 2
	}

	if rand.Intn(100) < enemy.Dodge {
		fmt.Println("Enemy Dodged!")
	} else {
		// Apply the damage to the enemy.
		fmt.Printf("Attacking %s for %d damage.\n", enemy.Name, damage)
		*enemyHp -= damage
	}
}

func enemyAttacks(player models.Player, enemy models.Enemy, playerHp *int, enemyHp *int) {
	// Calculate the enemy's damage.
	damage := enemy.Attack - player.Defense
	if rand.Intn(100) < enemy.Crit {
		damage *= 2
	}

	if rand.Intn(100) < player.Dodge {
		fmt.Println("Player Dodged!")
	} else {
		// Apply the damage to the player.
		fmt.Printf("Attacking %s for %d damage.\n", player.Name, damage)
		*playerHp -= damage
	}
}
