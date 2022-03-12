package tasks

import (
	"leet/models"
	"leet/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TasksRegister(group *gin.RouterGroup) {
	tasks := group.Group("/tasks")
	tasks.GET("/", GetTasks)
	tasks.POST("/", CreateTask)
	tasks.PUT("/:id", UpdateTask)
	tasks.DELETE("/:id", DeleteTask)

}

func GetTasks(c *gin.Context) {
	util.Log.Info().Msg("GetTasks " + c.ClientIP())

	var tasks []models.Task
	db := util.GetDB()
	db.Find(&tasks)
	c.JSON(200, tasks)
}

func CreateTask(c *gin.Context) {
	util.Log.Info().Msg("CreateTask " + c.ClientIP())
	var task models.Task
	var db = util.GetDB()

	if err := c.BindJSON(&task); err != nil {
		util.Log.Error().Msg("CreateTask 400 BindJSON " + c.ClientIP())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	db.Create(&task)
	c.JSON(http.StatusOK, &task)
}

func UpdateTask(c *gin.Context) {
	util.Log.Info().Msg("UpdateTask " + c.ClientIP())
	id := c.Param("id")
	var task models.Task

	db := util.GetDB()
	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		util.Log.Error().Msg("CreateTask 400 FindInDB " + c.ClientIP())
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.BindJSON(&task); err != nil {
		util.Log.Error().Msg("CreateTask 500 BindJSON from DB Obj " + c.ClientIP())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	task.UpdatedAt = time.Now()
	db.Save(&task)
	c.JSON(http.StatusOK, &task)
}

func DeleteTask(c *gin.Context) {
	util.Log.Info().Msg("DeleteTask " + c.ClientIP())
	id := c.Param("id")
	var task models.Task
	db := util.GetDB()

	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		util.Log.Error().Msg("DeleteTask 400 FindInDB " + c.ClientIP())
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	db.Delete(&task)
}
