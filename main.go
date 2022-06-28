package main

import (
	"fmt"
	"leet/tasks"
	"leet/util"
	"net/http"
	"os"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func health(c *gin.Context) {
	util.Log.Info().Msg(fmt.Sprintf("id:%s, health", requestid.Get(c)))
	c.JSON(http.StatusOK, gin.H{"message": "healthy"})
}

func main() {
	r := gin.Default()
	util.InitLogger(r)
	util.Log.Info().Msg("Starting server..")

	defer util.DBFree()
	connectionString := os.Getenv("POSTGRES_CONNECTION")
	fmt.Println("connectionString: ", connectionString)
	if connectionString == "" {
		postgresPassword := os.Getenv("POSTGRES_PASSWORD")
		connectionString = fmt.Sprintf("tasks-postgres-master.tasks.svc.cluster.local port=5432 user=app-user dbname=tasks password=%s sslmode=disable", postgresPassword)
	}
	util.InitPostgres(connectionString)

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
