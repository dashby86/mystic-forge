package sql

import (
	"database/sql"
	"fmt"
	"log"
	"mf/models"
)

type SqlService struct {
	DB *sql.DB
}

func (s SqlService) GetPlayerByID() (models.Player, error) {
	fmt.Println("Querying...")
	player := models.Player{}
	// Prepare the SQL statement
	stmt, err := s.DB.Prepare("SELECT id, name, player_level, forge_level FROM player WHERE id = 1")
	if err != nil {
		return player, fmt.Errorf("failed to prepare the SQL statement: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)
	// Execute the query
	err = stmt.QueryRow().Scan(&player.Id, &player.Name, &player.Level, &player.ForgeLevel)
	if err != nil {
		return player, fmt.Errorf("failed to execute the query: %v", err)
	}

	stmt, err = s.DB.Prepare("SELECT " +
		"SUM(hp), " +
		"SUM(attack), " +
		"SUM(defense), " +
		"SUM(speed), " +
		"SUM(crit), " +
		"SUM(dodge), " +
		"SUM(block) " +
		"FROM player_gear WHERE player_id = 1")

	if err != nil {
		return player, fmt.Errorf("failed to prepare the SQL statement: %v", err)
	}

	err = stmt.QueryRow().Scan(&player.HP, &player.Attack, &player.Defense, &player.Speed, &player.Crit, &player.Dodge, &player.Block)
	if err != nil {
		return player, fmt.Errorf("failed to execute the query: %v", err)
	}

	stmt, err = s.DB.Prepare("SELECT quantity FROM ore_inventory where player_id = 1")
	err = stmt.QueryRow().Scan(&player.Ore)
	if err != nil {
		return player, fmt.Errorf("failed to execute the query: %v", err)
	}
	return player, nil
}

func (s SqlService) SaveGearToSlot(player models.Player, gear models.Gear) (bool, error) {
	stmt, err := s.DB.Prepare("INSERT INTO player_gear (" +
		"player_id, " +
		"gear_slot_id, " +
		"rarity_id, " +
		"hp, " +
		"attack, " +
		"defense, " +
		"speed, " +
		"crit, " +
		"dodge, " +
		"block" +
		")" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE " +
		"rarity_id = ?, " +
		"hp = ?, " +
		"attack = ?, " +
		"defense = ?, " +
		"speed = ?, " +
		"crit = ?, " +
		"dodge = ?, " +
		"block = ? ")
	//"ON DUPLICATE KEY UPDATE ")
	if err != nil {
		return false, fmt.Errorf("failed to prepare the SQL statement: %v", err)
	}
	defer stmt.Close()
	// Execute the query
	_, err = stmt.Exec(
		player.Id,
		gear.SlotId,
		gear.Rarity,
		gear.HP,
		gear.Attack,
		gear.Defense,
		gear.Speed,
		gear.Crit,
		gear.Dodge,
		gear.Level,
		gear.Rarity,
		gear.HP,
		gear.Attack,
		gear.Defense,
		gear.Speed,
		gear.Crit,
		gear.Dodge,
		gear.Level)
	if err != nil {
		log.Fatalf("Failed to query the database: %v", err)
	}
	return true, nil
}

func (s SqlService) SpendOre(playerId int) (bool, error) {
	stmt, err := s.DB.Prepare("UPDATE ore_inventory SET quantity = quantity - 1 where player_id = ?")
	_, err = stmt.Exec(playerId)
	if err != nil {
		return false, fmt.Errorf("failed to execute the query: %v", err)
	}
	return true, nil
}

func (s SqlService) GrantOre(playerId int, amount int) (bool, error) {
	stmt, err := s.DB.Prepare("UPDATE ore_inventory SET quantity = quantity + ? where player_id = ?")
	_, err = stmt.Exec(amount, playerId)
	if err != nil {
		return false, fmt.Errorf("failed to execute the query: %v", err)
	}
	return true, nil
}

func (s SqlService) GetEquipedGearBySlot(playerId int, slotId int) (models.Gear, error) {
	gear := models.Gear{}
	stmt, err := s.DB.Prepare("SELECT " +
		"hp, " +
		"attack, " +
		"defense, " +
		"speed, " +
		"crit, " +
		"dodge, " +
		"block " +
		"FROM player_gear WHERE player_id = ? AND gear_slot_id = ?")

	if err != nil {
		return gear, fmt.Errorf("failed to prepare the SQL statement: %v", err)
	}
	err = stmt.QueryRow(playerId, slotId).Scan(&gear.HP, &gear.Attack, &gear.Defense, &gear.Speed, &gear.Crit, &gear.Dodge, &gear.Block)
	if err != nil {
		return gear, fmt.Errorf("failed to execute the query: %v", err)
	}
	return gear, nil
}

func (s SqlService) GrantExp(playerId int, exp int) error {
	stmt, err := s.DB.Prepare("UPDATE PLAYER SET player_exp = player_exp + ? WHERE id = ?")
	_, err = stmt.Exec(exp, playerId)
	if err != nil {
		return fmt.Errorf("failed to execute the query: %v", err)
	}
	return nil
}
