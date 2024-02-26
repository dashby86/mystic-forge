package sql

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"mf/models"
)

type SqlService struct {
	DB *sql.DB
}

func (s SqlService) GetPlayerByID() (models.Player, error) {
	fmt.Println("Querying...")
	player := models.Player{}
	// Prepare the SQL statement
	stmt, err := s.DB.Prepare("SELECT id, name, dungeon_level, player_level, forge_level, player_exp FROM player WHERE id = 1")
	if err != nil {
		return player, fmt.Errorf("failed to prepare the SQL statement: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)
	// Execute the query
	err = stmt.QueryRow().Scan(&player.Id, &player.Name, &player.DungeonLevel, &player.Level, &player.ForgeLevel, &player.Experience)
	if err != nil {
		return player, fmt.Errorf("failed to execute the query: %v", err)
	}

	stmt, err = s.DB.Prepare("SELECT " +
		"COALESCE(SUM(hp), 0), " +
		"COALESCE(SUM(attack), 0), " +
		"COALESCE(SUM(defense), 0), " +
		"COALESCE(SUM(speed), 0), " +
		"COALESCE(CAST(SUM(crit) AS DECIMAL(5,2)), 0), " +
		"COALESCE(CAST(SUM(dodge) AS DECIMAL(5,2)), 0), " +
		"COALESCE(CAST(SUM(block) AS DECIMAL(5,2)), 0) " +
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
		"player_id, gear_slot_id, level, rarity_id, hp, attack, defense, speed, crit, dodge, block" +
		") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE " +
		"player_id = ?, gear_slot_id = ?, level = ?, rarity_id = ?, hp = ?, attack = ?, defense = ?, speed = ?, crit = ?, dodge = ?, block = ?")

	if err != nil {
		log.Fatal("failed to prepare the SQL statement: %v", err)
	}
	defer stmt.Close() // Ensure statement closure

	_, err = stmt.Exec(
		player.Id, gear.SlotId, gear.Level, gear.Rarity, gear.HP, gear.Attack, gear.Defense, gear.Speed, gear.Crit, gear.Dodge, gear.Block,
		player.Id, gear.SlotId, gear.Level, gear.Rarity, gear.HP, gear.Attack, gear.Defense, gear.Speed, gear.Crit, gear.Dodge, gear.Block,
	)

	if err != nil {
		log.Fatalf("Failed to query the database: %v", err)
		return false, err
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
		"COALESCE(hp, 0), " +
		"COALESCE(level, 0), " +
		"COALESCE(rarity_id, 0), " +
		"COALESCE(attack, 0), " +
		"COALESCE(defense, 0), " +
		"COALESCE(speed, 0), " +
		"COALESCE(crit, 0), " +
		"COALESCE(dodge, 0), " +
		"COALESCE(block, 0) " +
		"FROM player_gear WHERE player_id = ? AND gear_slot_id = ?")

	if err != nil {
		return gear, fmt.Errorf("failed to prepare the SQL statement: %v", err)
	}
	err = stmt.QueryRow(playerId, slotId).Scan(&gear.HP, &gear.Level, &gear.Rarity, &gear.Attack, &gear.Defense, &gear.Speed, &gear.Crit, &gear.Dodge, &gear.Block)
	if err != nil && err != sql.ErrNoRows {
		return gear, fmt.Errorf("failed to execute the query: %v", err)
	}
	return gear, nil
}

func (s SqlService) GrantExp(playerId int, exp int) error {
	fmt.Println("UPDATE player SET player_exp = player_exp + ", +exp, " WHERE id = ", playerId)
	stmt, err := s.DB.Prepare("UPDATE player SET player_exp = player_exp + ? WHERE id = ?")
	_, err = stmt.Exec(exp, playerId)
	if err != nil {
		return fmt.Errorf("failed to execute the query: %v", err)
	}
	levelUpPlayer(s, playerId)
	return nil
}
func (s SqlService) GrantForgeExp(playerId int, exp int) error {
	//stmt, err := s.DB.Prepare("UPDATE PLAYER SET player_exp = player_exp + ? WHERE id = ?")
	//_, err = stmt.Exec(exp, playerId)
	//if err != nil {
	//	return fmt.Errorf("failed to execute the query: %v", err)
	//}
	return nil
}

func calculateRequiredExp(targetLevel int) int {
	exponent := 1.1
	baseExp := 20

	requiredExp := float64(baseExp) * math.Pow(float64(targetLevel), exponent)
	return int(requiredExp)
}

func levelUpPlayer(s SqlService, playerID int) (bool, error) {
	// 1. Begin Transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return false, err
	}

	// 2. Retrieve Player (within transaction)
	var currentLevel, currentExp int
	err = tx.QueryRow("SELECT player_level, player_exp FROM player WHERE id = ? FOR UPDATE", playerID).Scan(&currentLevel, &currentExp)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return false, err
	}

	// 3. Calculate Required Exp and Level Up
	requiredExp := calculateRequiredExp(currentLevel + 1)
	leveledUp := false // Flag to track if a level up occurred

	for currentExp >= requiredExp {
		currentLevel++
		currentExp -= requiredExp
		requiredExp = calculateRequiredExp(currentLevel + 1)
		leveledUp = true
	}

	// 4. Update Database within transaction
	if leveledUp {
		_, err = tx.Exec("UPDATE player SET player_level = ?, player_exp = ? WHERE id = ?", currentLevel, currentExp, playerID)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return false, err
		}
	}

	// 5. Commit Transaction
	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return leveledUp, nil
}

func (s SqlService) UpdateDungeonLevel(playerID int, dungeonLevel int) error {
	stmt, err := s.DB.Prepare("UPDATE player SET dungeon_level = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare the SQL statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(dungeonLevel, playerID)
	if err != nil {
		return fmt.Errorf("failed to execute the query: %v", err)
	}
	return nil
}
