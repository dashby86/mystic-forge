package mfsql

import (
	"database/sql"
	"fmt"
	"mf/models"
)

const (
	DBUser     = "app"       // Database username
	DBPassword = "admin"     // Database password
	DBHost     = "127.0.0.1" // Docker container hostname
	DBPort     = "3309"      // MySQL port
	DBName     = "app"       // Name of the database
)

type SqlService struct {
	DB              *sql.DB
	getPlayer       *sql.Stmt
	getPlayerHP     *sql.Stmt
	getOre          *sql.Stmt
	saveGearToSlot  *sql.Stmt
	spendOre        *sql.Stmt
	grantOre        *sql.Stmt
	getEquippedGear *sql.Stmt
}

func NewDBConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUser, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	return db, nil
}

func NewSqlService(db *sql.DB) (*SqlService, error) {
	s := &SqlService{DB: db}
	if err := s.initStatements(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SqlService) initStatements() error {
	var err error

	s.getPlayer, err = s.DB.Prepare(GET_PLAYER_BY_ID)
	if err != nil {
		return fmt.Errorf("failed to prepare getPlayer statement: %v", err)
	}

	s.getPlayerHP, err = s.DB.Prepare(GET_PLAYER_GEAR_BY_ID)
	if err != nil {
		return fmt.Errorf("failed to prepare getPlayerHP statement: %v", err)
	}

	s.getOre, err = s.DB.Prepare(GET_ORE_QUANTITY)
	if err != nil {
		return fmt.Errorf("failed to prepare getOre statement: %v", err)
	}

	s.saveGearToSlot, err = s.DB.Prepare(SAVE_GEAR_TO_SLOT)
	if err != nil {
		return fmt.Errorf("failed to prepare saveGearToSlot statement: %v", err)
	}

	s.spendOre, err = s.DB.Prepare(DECREMENT_ORE_QUANTITY)
	if err != nil {
		return fmt.Errorf("failed to prepare spendOre statement: %v", err)
	}

	s.grantOre, err = s.DB.Prepare(UPDATE_ORE_QUANTITY)
	if err != nil {
		return fmt.Errorf("failed to prepare grantOre statement: %v", err)
	}

	s.getEquippedGear, err = s.DB.Prepare(GET_EQUIPPED_GEAR)
	if err != nil {
		return fmt.Errorf("failed to prepare getEquippedGear statement: %v", err)
	}

	return nil
}

func (s *SqlService) Close() {
	if s.getPlayer != nil {
		_ = s.getPlayer.Close()
	}
	if s.getPlayerHP != nil {
		_ = s.getPlayerHP.Close()
	}
	if s.getOre != nil {
		_ = s.getOre.Close()
	}
	if s.saveGearToSlot != nil {
		_ = s.saveGearToSlot.Close()
	}
	if s.spendOre != nil {
		_ = s.spendOre.Close()
	}
	if s.grantOre != nil {
		_ = s.grantOre.Close()
	}
	if s.getEquippedGear != nil {
		_ = s.getEquippedGear.Close()
	}
}

func (s *SqlService) GetPlayerByID(playerID int) (models.Player, error) {
	player := models.Player{}

	err := s.getPlayer.QueryRow(playerID).Scan(&player.Id, &player.Name)
	if err != nil {
		return player, fmt.Errorf("failed to find player by id: %v", err)
	}

	err = s.getPlayerHP.QueryRow(playerID).Scan(&player.MaxHP, &player.Attack, &player.Defense, &player.Speed, &player.Crit, &player.Dodge, &player.Block)
	if err != nil {
		return player, fmt.Errorf("failed to execute getPlayerHP query: %v", err)
	}

	err = s.getOre.QueryRow(playerID).Scan(&player.Ore)
	if err != nil && err != sql.ErrNoRows {
		return player, fmt.Errorf("failed to execute getOre query: %v", err)
	}

	return player, nil
}

func (s SqlService) SaveGearToSlot(player models.Player, gear models.Gear) (bool, error) {
	if s.saveGearToSlot == nil {
		if err := s.initStatements(); err != nil {
			return false, err
		}
	}

	_, err := s.saveGearToSlot.Exec(
		player.Id,
		gear.SlotId,
		gear.Rarity,
		gear.HP,
		gear.Attack,
		gear.Defense,
		gear.Speed,
		gear.Crit,
		gear.Dodge,
		gear.Block,
		gear.Rarity,
		gear.HP,
		gear.Attack,
		gear.Defense,
		gear.Speed,
		gear.Crit,
		gear.Dodge,
		gear.Block,
	)
	if err != nil {
		return false, fmt.Errorf("failed to execute SaveGearToSlot query: %v", err)
	}
	return true, nil
}

func (s SqlService) SpendOre(playerId int) (bool, error) {
	if s.spendOre == nil {
		if err := s.initStatements(); err != nil {
			return false, err
		}
	}

	_, err := s.spendOre.Exec(playerId)
	if err != nil {
		return false, fmt.Errorf("failed to execute SpendOre query: %v", err)
	}
	return true, nil
}

func (s SqlService) GrantOre(playerId int, amount int) (bool, error) {
	if s.grantOre == nil {
		if err := s.initStatements(); err != nil {
			return false, err
		}
	}

	_, err := s.grantOre.Exec(amount, playerId)
	if err != nil {
		return false, fmt.Errorf("failed to execute GrantOre query: %v", err)
	}
	return true, nil
}

func (s SqlService) GetEquippedGearBySlot(playerId int, slotId int) (models.Gear, error) {
	gear := models.Gear{}
	if s.getEquippedGear == nil {
		if err := s.initStatements(); err != nil {
			return gear, err
		}
	}

	err := s.getEquippedGear.QueryRow(playerId, slotId).Scan(&gear.HP, &gear.Attack, &gear.Defense, &gear.Speed, &gear.Crit, &gear.Dodge, &gear.Block)
	if err != nil {
		return gear, fmt.Errorf("failed to execute GetEquippedGearBySlot query: %v", err)
	}
	return gear, nil
}
