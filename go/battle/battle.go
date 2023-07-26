package battle

import (
	"fmt"
	"math/rand"
	"mf/models"
)

type Battle struct {
	Player models.Player
	Enemey models.Enemy
}

func (battle Battle) SimBattle() {
	// Initialize the turn variable.
	turn := 1

	// Simulate the battle.
	for {
		// Check which entity gets to attack first.
		if battle.Player.Speed > battle.Enemey.Speed {
			playerAttacks(battle.Player, battle.Enemey)
		} else {
			enemyAttacks(battle.Player, battle.Enemey)
		}

		// Check if the battle is over.
		if battle.Player.HP <= 0 {
			fmt.Println("The player has died.")
			break
		}
		if battle.Enemey.HP <= 0 {
			fmt.Println("The enemy has died.")
			break
		}

		// Switch turns.
		turn = 3 - turn
	}
}

func playerAttacks(player models.Player, enemy models.Enemy) {
	// Calculate the player's damage.
	damage := player.Attack - enemy.Defense
	if rand.Intn(100) < player.Crit {
		fmt.Println("Critical Strike!")
		damage *= 2
	}

	// Apply the damage to the enemy.
	fmt.Printf("Attacking for %s damage.\n", damage)
	enemy.HP -= damage
}

func enemyAttacks(player models.Player, enemy models.Enemy) {
	// Calculate the enemy's damage.
	damage := enemy.Attack - player.Defense
	if rand.Intn(100) < enemy.Crit {
		damage *= 2
	}

	// Apply the damage to the player.
	fmt.Printf("Attacking for %s damage.\n", damage)
	player.HP -= damage
}
