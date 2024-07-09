package database

import (
	"context"
	"database/sql"
	"os"
	"path"
	"time"

	"github.com/danilotorchio/goexpert-clientserverapi/models"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	dbFile := path.Join(dir, "database", "database.db")

	_, err = os.Stat(dbFile)
	if os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			return err
		}
		file.Close()
	}

	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	err = createTables()
	return err
}

func createTables() error {
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS hist_usd_brl (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bid TEXT,
			timestamp INTEGER
		);
	`

	_, err := DB.Exec(sqlStmt)
	return err
}

func InsertNewExchangeRate(ctx context.Context, bid string) error {
	stmt, err := DB.PrepareContext(ctx, "INSERT INTO hist_usd_brl(bid, timestamp) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(bid, time.Now().Unix())
	return err
}

func GetExchangeHistory(ctx context.Context, limit int) (*[]models.HistUsdBrl, error) {
	limit = max(limit, 1)
	limit = min(limit, 100)

	stmt, err := DB.PrepareContext(ctx, "SELECT bid, timestamp FROM hist_usd_brl ORDER BY timestamp DESC LIMIT ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.HistUsdBrl

	for rows.Next() {
		var d models.HistUsdBrl
		if err := rows.Scan(&d.Bid, &d.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, d)
	}

	return &history, nil
}
