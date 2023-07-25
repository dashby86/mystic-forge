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
	player := models.Player{}
	// Prepare the SQL statement
	stmt, err := s.DB.Prepare("SELECT id, name FROM player WHERE id = 1")
	if err != nil {
		return player, fmt.Errorf("failed to prepare the SQL statement: %v", err)
	}
	defer stmt.Close()
	// Execute the query
	err = stmt.QueryRow().Scan(&player.Id, &player.Name)
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
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ")
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
		gear.Level)
	if err != nil {
		log.Fatalf("Failed to query the database: %v", err)
	}
	return true, nil
}
