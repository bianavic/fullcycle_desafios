package config

import (
	"fmt"
	"github.com/bianavic/fullcycle_desafios.git/schemas"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

const (
	dbSourceName = "./Client-Server-API/data/db/exchange_rates.db?cache=shared&mode=memory" // sqlite db file
)

func InitializeSQLite() (*gorm.DB, error) {
	logger := GetLogger("sqlite")

	// ensure file exists
	_, err := os.Stat(dbSourceName)
	if os.IsNotExist(err) {
		logger.Info("databse file not found, creating...")
		// create db
		err := os.MkdirAll("./Client-Server-API/data/db", os.ModePerm)
		if err != nil {
			return nil, err
		}
		file, err := os.Create(dbSourceName)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	// open sqlite connection
	db, err := gorm.Open(sqlite.Open(dbSourceName), &gorm.Config{})
	if err != nil {
		logger.Errorf("error opening sqlite %v", err)
		return nil, err
	}

	// check tables
	if !db.Migrator().HasTable(&schemas.Rate{}) {
		logger.Info("tables not found, creating...")
		err := db.Migrator().CreateTable(&schemas.Rate{})
		if err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return db, nil
}
