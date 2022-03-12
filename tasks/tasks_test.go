package tasks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"leet/db"
	"leet/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTaskModel(t *testing.T) {
	// get test db
	testDB := db.InitTestDB()
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
	db.DBFree(testDB)

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
		"{}",
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
			asserts.Equal(status, http.StatusOK)
			task, err := ParseTask(body)
			asserts.NoError(err)
			asserts.Equal("test", task.Title)
		},
	},
	{
		"get task just created from db",
		"GET",
		"/api/v1/tasks/",
		"{}",
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
		"/api/v1/tasks/1",
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
		`{}`,
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
		"/api/v1/tasks/1",
		`{}`,
		func(t *testing.T, status int, body string) {
			asserts := assert.New(t)
			asserts.Equal(status, http.StatusOK)
		},
	},
	{
		"verify no tasks left",
		"GET",
		"/api/v1/tasks/",
		`{}`,
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
	testDB := db.InitTestDB()
	asserts := assert.New(t)

	// migrate
	testDB.AutoMigrate(&models.Task{})

	// setup router
	r := gin.New()
	v1 := r.Group("/api/v1")
	TasksRegister(v1)

	for _, workflow := range TaskRequestWorkflow {
		request, err := http.NewRequest(workflow.method, workflow.url,
			bytes.NewBufferString(workflow.body))
		request.Header.Set("Content-Type", "application/json")
		asserts.NoError(err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		workflow.assertCorrect(t, w.Code, w.Body.String())
	}

	// close db
	db.DBFree(testDB)

}
