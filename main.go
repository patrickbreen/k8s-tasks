package main

import (
	"leet/tasks"
	"leet/util"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func health(c *gin.Context) {
	util.Log.Info().Msg("Checked health")
	c.JSON(http.StatusOK, gin.H{"message": "healthy"})
}

func main() {
	util.Log = util.InitLogger()
	util.Log.Info().Msg("Starting server..")

	defer util.DBFree()
	connectionString := os.Getenv("POSTGRES_URL")
	util.InitPostgres(connectionString)

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
	tasks.TasksRegister(v1)

	if err := r.Run(); err != nil {
		util.Log.Fatal().Msg("main error: " + err.Error())
	}
}
