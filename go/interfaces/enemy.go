package interfaces

type Enemy interface {
	GetId() int        // Method to get the enemy's id.
	GetName() string   // Method to get the name of the enemy.
	GetCurrentHP() int // Method to get the current HP of the enemy.
	GetMaxHP() int     // Method to get the maximum HP of the enemy.
	GetAttack() int    // Method to get the enemy's attack value.
	GetDefense() int   // Method to get the enemy's defense value.
	GetSpeed() float64 // Method to get the enemy's attack speed value.
	GetCrit() float64  // Method to get the enemy's critical chance value.
	GetDodge() float64 // Method to get the enemy's dodge chance value.
	GetBlock() float64 // Method to get the enemy's block chance value.

	Defend(damage int) // Method to defend against an attack.
}
