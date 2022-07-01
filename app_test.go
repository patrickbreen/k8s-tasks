package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

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

	// create
	request, err := http.NewRequest("POST",
		serverDomain+"/api/v1/tasks/",
		bytes.NewBufferString(`{"Title": "test", "Completed": false}`))
	request.Header.Set("Content-Type", "application/json")
	asserts.NoError(err)
	c := &http.Client{}
	response, err := c.Do(request)
	asserts.NoError(err)
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	asserts.Equal(http.StatusOK, response.StatusCode)
	var task models.Task
	err = json.Unmarshal(buf.Bytes(), &task)
	asserts.NoError(err)
	asserts.Equal("test", task.Title)

	// verify get, TODO this should be a lookup by ID
	request, err = http.NewRequest("GET",
		serverDomain+"/api/v1/tasks/",
		bytes.NewBufferString(``))
	request.Header.Set("Content-Type", "application/json")
	asserts.NoError(err)
	response, err = c.Do(request)
	asserts.NoError(err)
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	asserts.Equal(http.StatusOK, response.StatusCode)
	var tasks []models.Task
	err = json.Unmarshal(buf.Bytes(), &tasks)
	asserts.NoError(err)
	foundTask := false
	for _, returnedTask := range tasks {
		if returnedTask.ID == task.ID {
			foundTask = true
		}
	}
	asserts.Equal(true, foundTask)

	// update
	// id := task.ID
	request, err = http.NewRequest("PUT",
		serverDomain+"/api/v1/tasks/?id="+fmt.Sprint(task.ID),
		bytes.NewBufferString(`{"title": "changedit", "completed": false}`))
	request.Header.Set("Content-Type", "application/json")
	asserts.NoError(err)
	response, err = c.Do(request)
	asserts.NoError(err)
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	asserts.Equal(http.StatusOK, response.StatusCode)
	err = json.Unmarshal(buf.Bytes(), &task)
	asserts.NoError(err)
	asserts.Equal("changedit", task.Title)

	// delete
	request, err = http.NewRequest("DELETE",
		serverDomain+"/api/v1/tasks/?id="+fmt.Sprint(task.ID),
		bytes.NewBufferString(``))
	request.Header.Set("Content-Type", "application/json")
	asserts.NoError(err)
	response, err = c.Do(request)
	asserts.NoError(err)
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	asserts.Equal(http.StatusOK, response.StatusCode)

	// verify no get, TODO this should be a lookup by ID
	request, err = http.NewRequest("GET",
		serverDomain+"/api/v1/tasks/",
		bytes.NewBufferString(``))
	request.Header.Set("Content-Type", "application/json")
	asserts.NoError(err)
	response, err = c.Do(request)
	asserts.NoError(err)
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	asserts.Equal(http.StatusOK, response.StatusCode)
	err = json.Unmarshal(buf.Bytes(), &tasks)
	asserts.NoError(err)
	foundTask = false
	for _, returnedTask := range tasks {
		if returnedTask.ID == task.ID {
			foundTask = true
		}
	}
	asserts.Equal(false, foundTask)
}
