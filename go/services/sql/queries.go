package mfsql

const (
	GET_PLAYER_BY_ID      = "SELECT id, name FROM player WHERE id = ?"
	GET_PLAYER_GEAR_BY_ID = `SELECT
	COALESCE(SUM(hp), 0) AS sum_hp,
	COALESCE(SUM(attack), 0) AS sum_attack,
	COALESCE(SUM(defense), 0) AS sum_defense,
	COALESCE(SUM(speed), 0) AS sum_speed,
	COALESCE(SUM(crit), 0) AS sum_crit,
	COALESCE(SUM(dodge), 0) AS sum_dodge,
	COALESCE(SUM(block), 0) AS sum_block
FROM
	player_gear
WHERE
	player_id = ?`
	GET_ORE_QUANTITY       = "SELECT quantity FROM ore_inventory WHERE player_id = ?"
	SAVE_GEAR_TO_SLOT      = "INSERT INTO player_gear (player_id, gear_slot_id, rarity_id, hp, attack, defense, speed, crit, dodge, block) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE rarity_id = ?, hp = ?, attack = ?, defense = ?, speed = ?, crit = ?, dodge = ?, block = ?"
	DECREMENT_ORE_QUANTITY = "UPDATE ore_inventory SET quantity = quantity - 1 WHERE player_id = ?"
	UPDATE_ORE_QUANTITY    = "UPDATE ore_inventory SET quantity = quantity + ? WHERE player_id = ?"
	GET_EQUIPPED_GEAR      = "SELECT hp, attack, defense, speed, crit, dodge, block FROM player_gear WHERE player_id = ? AND gear_slot_id = ?"
)
