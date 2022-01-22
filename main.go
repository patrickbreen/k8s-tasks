package main

import (
	Controller "leet/controllers"
	"leet/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func main() {
	log.Println("Starting server..")

	db.Init()

	r := gin.Default()
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(r)

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
