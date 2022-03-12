package tasks

import (
	"fmt"
	"leet/models"
	"leet/util"
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
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
	util.Log.Info().Msg(fmt.Sprintf("%s, GetTasks", requestid.Get(c)))

	var tasks []models.Task
	db := util.GetDB()
	db.Find(&tasks)
	c.JSON(200, tasks)
}

func CreateTask(c *gin.Context) {
	util.Log.Info().Msg(fmt.Sprintf("%s, CreateTask", requestid.Get(c)))
	var task models.Task
	var db = util.GetDB()

	if err := c.BindJSON(&task); err != nil {
		util.Log.Error().Msg(fmt.Sprintf("%s, CreateTask 400 BindJSON", requestid.Get(c)))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	db.Create(&task)
	c.JSON(http.StatusOK, &task)
}

func UpdateTask(c *gin.Context) {
	util.Log.Info().Msg(fmt.Sprintf("%s, UpdateTask", requestid.Get(c)))
	id := c.Param("id")
	var task models.Task

	db := util.GetDB()
	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		util.Log.Error().Msg(fmt.Sprintf("%s, UpdateTask 403 DB lookup", requestid.Get(c)))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.BindJSON(&task); err != nil {
		util.Log.Error().Msg(fmt.Sprintf("%s, UpdateTask 500 BindJSON", requestid.Get(c)))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	task.UpdatedAt = time.Now()
	db.Save(&task)
	c.JSON(http.StatusOK, &task)
}

func DeleteTask(c *gin.Context) {
	util.Log.Info().Msg(fmt.Sprintf("%s, DeleteTask", requestid.Get(c)))
	id := c.Param("id")
	var task models.Task
	db := util.GetDB()

	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		util.Log.Error().Msg(fmt.Sprintf("%s, DeleteTask 403 DB lookup", requestid.Get(c)))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	db.Delete(&task)
}
