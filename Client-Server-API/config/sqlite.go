package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

const (
	dbSourceName = "./Client-Server-API/data/db/exchange_rates.db?cache=shared&mode=memory" // sqlite db file
)

func InitializeSQLite() (*gorm.DB, error) {
	logger := GetLogger("sqlite")

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

	// create db and connect
	db, err := gorm.Open(sqlite.Open(dbSourceName), &gorm.Config{})
	if err != nil {
		logger.Errorf("error opening sqlite %v", err)
		return nil, err
	}

	// migrate schema
	err = db.AutoMigrate()
	if err != nil {
		logger.Errorf("sqlite automigration error: %v", err)
		return nil, err
	}

	return db, nil
}
