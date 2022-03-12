package util

import (
	"leet/models"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// globals hold the db reference.
var db *gorm.DB
var err error

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// creates a connection to postgres database and migrates any new models
func InitPostgres(connectionString string) {

	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Log.Info().Msg("Database connected")

	if err := db.AutoMigrate(&models.Task{}); err != nil {
		Log.Info().Msg(err.Error())
	}

	// set connection limits
	if dbObj, err := db.DB(); err == nil {
		dbObj.SetMaxIdleConns(10)
		dbObj.SetMaxOpenConns(100)
	}

}

// sqlite is for testing or when you don't want to run postgress
func InitTestDB() *gorm.DB {
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if dbObj, err := db.DB(); err == nil {
		dbObj.SetMaxIdleConns(10)
		//dbObj.LogMode(true)
	}
	return db
}

func DBFree() error {
	sqlDB, err := db.DB()
	sqlDB.Close()
	if err != nil {
		panic(err)
	}
	return err
}

// GetDB object reference
func GetDB() *gorm.DB {
	return db
}
