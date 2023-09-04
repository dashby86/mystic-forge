package battle

import (
	"fmt"
	"math/rand"
	"mf/interfaces"
	"time"
)

type Battle struct {
	Player interfaces.Player
	Enemy  interfaces.Enemy
	Winner interfaces.Character
}

func (battle *Battle) SimBattle(eventQueue interfaces.EventQueue) {
	// Initialize entities and attack intervals.
	enemyHp := battle.Enemy.GetMaxHP()
	playerHp := battle.Player.GetMaxHP()
	playerAttackInterval := 1.0 / battle.Player.GetSpeed()
	enemyAttackInterval := 1.0 / battle.Enemy.GetSpeed()

	fmt.Printf("Starting battle! Player life: %d -- Enemy life: %d\n", playerHp, enemyHp)

	startTime := time.Now()

	for {
		elapsedTime := time.Since(startTime).Seconds()

		// Check if the player can attack.
		if elapsedTime >= playerAttackInterval {
			battle.playerAttacks()
			fmt.Printf("Player attacks! Enemy life: %d\n", battle.Enemy.GetCurrentHP())

			// Check if the enemy has been defeated.
			if battle.Enemy.GetCurrentHP() <= 0 {

				// TODO: Event - Enemy death
				fmt.Println("The enemy has died.")
				break
			}

			// Update the last attack time for the player.
			startTime = time.Now()
		}

		// Check if the enemy can attack.
		if elapsedTime >= enemyAttackInterval {
			battle.enemyAttacks()
			fmt.Printf("Enemy attacks! Player life: %d\n", battle.Player.GetCurrentHP())

			// Check if the player has been defeated.
			if battle.Player.GetCurrentHP() <= 0 {
				// TODO: Event - Player death
				fmt.Println("The player has died.")
				break
			}

			// Update the last attack time for the enemy.
			startTime = time.Now()
		}

		// Sleep briefly to avoid busy waiting and reduce resource usage.
		time.Sleep(100 * time.Millisecond)
	}
}

func (battle *Battle) playerAttacks() {
	// Calculate the player's damage.
	damage := battle.Player.GetAttack() - battle.Enemy.GetDefense()
	if rand.Float64() < battle.Player.GetCrit() {
		// TODO: Event - Player crits
		fmt.Println("Critical Strike!")
		damage *= 2
	}

	if rand.Float64() < battle.Enemy.GetDodge() {
		// TODO: Event - Enemy dodges
		fmt.Println("Enemy Dodged!")
	} else {
		// TODO: Event - Enemy damaged
		// Apply the damage to the enemy.
		fmt.Printf("Attacking %s for %d damage.\n", battle.Enemy.GetName(), damage)
		battle.Enemy.Defend(damage)
	}
}

func (battle *Battle) enemyAttacks() {
	// Calculate the enemy's damage.
	damage := battle.Enemy.GetAttack() - battle.Player.GetDefense()
	if rand.Float64() < battle.Enemy.GetCrit() {
		damage *= 2
	}

	if rand.Float64() < battle.Player.GetDodge() {
		// TODO: Event - Player dodges
		fmt.Println("Player Dodged!")
	} else {
		// TODO: Event - Player damaged
		// Apply the damage to the player.
		fmt.Printf("Attacking %s for %d damage.\n", battle.Player.GetName(), damage)
		battle.Player.Defend(damage)
	}
}
