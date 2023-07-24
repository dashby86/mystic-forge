package sql

import (
	"database/sql"
	"fmt"
	"mf/models"
)

type SqlService struct {
	DB *sql.DB
}

func (s SqlService) GetPlayerByID() (models.Player, error) {
	player := models.Player{}
	// Prepare the SQL statement
	stmt, err := s.DB.Prepare("SELECT id, name FROM Player WHERE id = 1")
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
