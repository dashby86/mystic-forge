package interfaces

type Player interface {
	GetId() int        // Method to get the player's id.
	GetName() string   // Method to get the name of the player.
	GetCurrentHP() int // Method to get the current HP of the player.
	GetMaxHP() int     // Method to get the maximum HP of the player.
	GetAttack() int    // Method to get the player's attack value.
	GetDefense() int   // Method to get the player's defense value.
	GetSpeed() float64 // Method to get the player's attack speed value.
	GetCrit() float64  // Method to get the player's critical chance value.
	GetDodge() float64 // Method to get the player's dodge chance value.
	GetBlock() float64 // Method to get the player's block chance value.
	GetOre() int       // Method to get the player's current ore count.

	Defend(damage int) // Method to defend against an attack.
}
