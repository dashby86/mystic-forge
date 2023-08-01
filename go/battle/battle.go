package battle

import (
	"fmt"
	"math/rand"
	"mf/events"
	"mf/models"
	sqlService "mf/services/sql"
	"time"
)

type Battle struct {
	Player models.Player
	Enemy  models.Enemy
	Sql    sqlService.SqlService
}

func (battle *Battle) SimBattle(eventQueue *events.EventQueue) {
	// Initialize entities and attack intervals.
	enemyHp := battle.Enemy.MaxHP
	playerHp := battle.Player.MaxHP
	playerAttackInterval := 1.0 / battle.Player.Speed
	enemyAttackInterval := 1.0 / battle.Enemy.Speed

	// TODO: Event - Battle Start
	eventQueue.AddEvent(
		events.BattleStartEvent{
			Player: battle.Player,
			Enemy:  battle.Enemy,
		},
	)

	fmt.Printf("Starting battle! Player life: %d -- Enemy life: %d\n", playerHp, enemyHp)

	startTime := time.Now()

	for {
		elapsedTime := time.Since(startTime).Seconds()

		// Check if the player can attack.
		if elapsedTime >= playerAttackInterval {
			battle.playerAttacks()
			fmt.Printf("Player attacks! Enemy life: %d\n", battle.Enemy.CurrentHP)

			// Check if the enemy has been defeated.
			if battle.Enemy.CurrentHP <= 0 {

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
			fmt.Printf("Enemy attacks! Player life: %d\n", battle.Player.CurrentHP)

			// Check if the player has been defeated.
			if battle.Player.CurrentHP <= 0 {
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

func (battle Battle) playerAttacks() {
	// Calculate the player's damage.
	damage := battle.Player.Attack - battle.Enemy.Defense
	if rand.Float64() < battle.Player.Crit {
		// TODO: Event - Player crits
		fmt.Println("Critical Strike!")
		damage *= 2
	}

	if rand.Float64() < battle.Enemy.Dodge {
		// TODO: Event - Enemy dodges
		fmt.Println("Enemy Dodged!")
	} else {
		// TODO: Event - Enemy damaged
		// Apply the damage to the enemy.
		fmt.Printf("Attacking %s for %d damage.\n", battle.Enemy.Name, damage)
		battle.Enemy.CurrentHP -= damage
	}
}

func (battle Battle) enemyAttacks() {
	// Calculate the enemy's damage.
	damage := battle.Enemy.Attack - battle.Player.Defense
	if rand.Float64() < battle.Enemy.Crit {
		damage *= 2
	}

	if rand.Float64() < battle.Player.Dodge {
		// TODO: Event - Player dodges
		fmt.Println("Player Dodged!")
	} else {
		// TODO: Event - Player damaged
		// Apply the damage to the player.
		fmt.Printf("Attacking %s for %d damage.\n", battle.Player.Name, damage)
		battle.Player.CurrentHP -= damage
	}
}
