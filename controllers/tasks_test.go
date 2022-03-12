package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"leet/db"
	"leet/models"
)

func TestGetCreateUpdateDelete(t *testing.T) {
	asserts := assert.New(t)
	asserts.Equal(1, 1)
	// get test db
	testDB := db.InitTestDB()

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
