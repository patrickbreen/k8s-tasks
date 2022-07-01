package main

import (
	"bytes"
	"encoding/json"
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

// name, method, url, body, expectedStatus, expectedBody
var TaskRequestWorkflow = []struct {
	name          string
	method        string
	url           string
	body          string
	assertCorrect func(*testing.T, int, string)
}{
	{
		"check no tasks exist",
		"GET",
		"/api/v1/tasks/",
		"",
		func(t *testing.T, status int, body string) {
			asserts := assert.New(t)
			asserts.Equal(status, http.StatusOK)
			tasks, err := ParseTasks(body)
			asserts.NoError(err)
			asserts.Equal(0, len(tasks))
		},
	},
	{
		"create",
		"POST",
		"/api/v1/tasks/",
		`{"Title": "test", "Completed": false}`,
		func(t *testing.T, status int, body string) {
			asserts := assert.New(t)
			asserts.Equal(http.StatusOK, status)
			task, err := ParseTask(body)
			asserts.NoError(err)
			asserts.Equal("test", task.Title)
		},
	},
	{
		"get task just created from db",
		"GET",
		"/api/v1/tasks/",
		"",
		func(t *testing.T, status int, body string) {
			asserts := assert.New(t)
			asserts.Equal(status, http.StatusOK)
			tasks, err := ParseTasks(body)
			asserts.NoError(err)
			asserts.Equal(1, len(tasks))
			asserts.Equal("test", tasks[0].Title)
		},
	},
	{
		"update",
		"PUT",
		"/api/v1/tasks/?id=1",
		`{"title": "changedit", "completed": false}`,
		func(t *testing.T, status int, body string) {
			asserts := assert.New(t)
			asserts.Equal(status, http.StatusOK)
			task, err := ParseTask(body)
			asserts.NoError(err)
			asserts.Equal("changedit", task.Title)
		},
	},
	{
		"get tasks just updated",
		"GET",
		"/api/v1/tasks/",
		``,
		func(t *testing.T, status int, body string) {
			asserts := assert.New(t)
			asserts.Equal(status, http.StatusOK)
			tasks, err := ParseTasks(body)
			asserts.NoError(err)
			asserts.Equal(1, len(tasks))
			asserts.Equal("changedit", tasks[0].Title)
		},
	},
	{
		"delete task",
		"DELETE",
		"/api/v1/tasks/?id=1",
		``,
		func(t *testing.T, status int, body string) {
			asserts := assert.New(t)
			asserts.Equal(status, http.StatusOK)
		},
	},
	{
		"verify no tasks left",
		"GET",
		"/api/v1/tasks/",
		``,
		func(t *testing.T, status int, body string) {
			asserts := assert.New(t)
			asserts.Equal(status, http.StatusOK)
			tasks, err := ParseTasks(body)
			asserts.NoError(err)
			asserts.Equal(0, len(tasks))
		},
	},
}

func ParseTask(s string) (models.Task, error) {
	var task models.Task
	err := json.Unmarshal([]byte(s), &task)
	return task, err
}

func ParseTasks(s string) ([]models.Task, error) {
	var tasks []models.Task
	err := json.Unmarshal([]byte(s), &tasks)
	return tasks, err
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

	for _, workflow := range TaskRequestWorkflow {
		request, err := http.NewRequest(workflow.method,
			"http://localhost:8080"+workflow.url,
			bytes.NewBufferString(workflow.body))
		request.Header.Set("Content-Type", "application/json")
		asserts.NoError(err)

		c := &http.Client{}
		response, err := c.Do(request)
		asserts.NoError(err)
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		workflow.assertCorrect(t, response.StatusCode, buf.String())
	}
}
