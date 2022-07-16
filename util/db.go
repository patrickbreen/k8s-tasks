package util

import (
	"fmt"
	"leet/models"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
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
func InitPostgres() {

	connectionString := os.Getenv("POSTGRES_CONNECTION")
	if connectionString == "" {
		postgresPassword := os.Getenv("POSTGRES_PASSWORD")
		connectionString = fmt.Sprintf("host=tasks-postgres-master.tasks.svc.cluster.local port=5432 user=owner dbname=app password=%s sslmode=disable", postgresPassword)
	}

	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Info().Msg("Database connected")

	if err := db.AutoMigrate(&models.Task{}); err != nil {
		log.Info().Msg(err.Error())
	}

	// set connection limits
	if dbObj, err := db.DB(); err == nil {
		dbObj.SetMaxIdleConns(5)
		dbObj.SetMaxOpenConns(20)
	}

}

// sqlite is for testing or when you don't want to run postgress
func InitTestDB() *gorm.DB {
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if dbObj, err := db.DB(); err == nil {
		dbObj.SetMaxIdleConns(5)
		//dbObj.LogMode(true)
	}
	return db
}

func DBFree() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Info().Msg("Error getting DB")
		panic(err)
	}
	err = sqlDB.Close()
	if err != nil {
		log.Info().Msg("Error closing DB")
		panic(err)
	}
}

// GetDB object reference
func GetDB() *gorm.DB {
	return db
}
