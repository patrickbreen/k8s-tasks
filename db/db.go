package db

import (
	"leet/models"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
		log.Println("db err: (InitPostgres) ", err)
		panic(err)
	}
	log.Println("Database connected")

	if err := db.AutoMigrate(&models.Task{}); err != nil {
		log.Println(err.Error())
	}

	// set connection limits
	if dbObj, err := db.DB(); err == nil {
		dbObj.SetMaxIdleConns(10)
		dbObj.SetMaxOpenConns(100)
	}

}

// sqlite is for testing or when you don't want to run postgress
func InitTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Println("db err: (TestDBInit) ", err)
	}
	if dbObj, err := db.DB(); err == nil {
		dbObj.SetMaxIdleConns(10)
		//dbObj.LogMode(true)
	}
	return db
}

func DBFree(db *gorm.DB) error {
	sqlDB, err := db.DB()
	sqlDB.Close()
	if err != nil {
		log.Println("db err: (DBFree) ", err)
	}
	return err
}

// GetDB object reference
func GetDB() *gorm.DB {
	return db
}
