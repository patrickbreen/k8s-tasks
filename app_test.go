package main

import (
	"net/http"
	"testing"
	"time"

	"leet/canary"
	"leet/models"
	"leet/util"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestTaskModel(t *testing.T) {
	// get test db
	testDB := util.InitTestDB()
	asserts := assert.New(t)

	// migrate
	testDB.AutoMigrate(&models.Task{})

	// check no tasks exist
	var tasks []models.Task
	testDB.Find(&tasks)
	asserts.Equal(0, len(tasks))

	// create
	task := &models.Task{Title: "test", Completed: false}
	task.CreatedAt = time.Now()
	testDB.Create(&task)

	// get task just created from DB
	testDB.Find(&tasks)
	asserts.Equal(1, len(tasks))
	asserts.Equal("test", tasks[0].Title)

	// update
	task.Title = "changedit"
	task.UpdatedAt = time.Now()
	testDB.Save(&task)

	// get task just updated from DB
	testDB.Find(&tasks)
	asserts.Equal(1, len(tasks))
	asserts.Equal("changedit", tasks[0].Title)

	// delete task
	testDB.Delete(&task)

	// check no tasks in db
	testDB.Find(&tasks)
	asserts.Equal(0, len(tasks))

	// close db
	util.DBFree()

}

func TestTaskRequests(t *testing.T) {
	// get test db - also just sets the global db reference
	testDB := util.InitTestDB()
	asserts := assert.New(t)

	// migrate
	testDB.AutoMigrate(&models.Task{})
	defer util.DBFree()

	// setup router
	InitPrometheus()
	mux := InitAppServer()
	log.Info().Msg("Server initialized")
	go http.ListenAndServe(":9000", promhttp.Handler())
	go http.ListenAndServe(":8080", mux)

	serverDomain := "http://localhost:8080"
	canary.RunCanary(serverDomain)
	asserts.Equal(true, true)
}
