package main

import (
	Controller "leet/controllers"
	"leet/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting server..")

	db.Init()

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.GET("/", Controller.GetTasks)
			tasks.POST("/", Controller.CreateTask)
			tasks.PUT("/:id", Controller.UpdateTask)
			tasks.DELETE("/:id", Controller.DeleteTask)
		}
	}

	r.Run()
}
