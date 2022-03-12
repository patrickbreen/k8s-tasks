package main

import (
	taskCtl "leet/controllers"
	"leet/db"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "healthy"})
}

func main() {
	log.Println("Starting server..")

	connectionString := os.Getenv("POSTGRES_URL")
	db.InitPostgres(connectionString)

	r := gin.Default()

	// gin metrics
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(r)

	r.GET("/health", health)

	r.Static("assets", "./assets")
	v1 := r.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.GET("/", taskCtl.GetTasks)
			tasks.POST("/", taskCtl.CreateTask)
			tasks.PUT("/:id", taskCtl.UpdateTask)
			tasks.DELETE("/:id", taskCtl.DeleteTask)
		}
	}

	if err := r.Run(); err != nil {
		log.Println("main error", err.Error())
	}
}
