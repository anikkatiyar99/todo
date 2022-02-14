package main

import (
	"fmt"
	"os"
	"time"

	"github.com/anikkatiyar99/todo/cron"
	routes "github.com/anikkatiyar99/todo/routes"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.TaskRoutes(router)

	// scheduler starts running jobs and current thread continues to execute
	fmt.Println("Cron Email Alert Running........")
	s1 := gocron.NewScheduler(time.UTC)
	s1.Every(50).Seconds().Do(cron.CronEmailRunner)
	s1.StartAsync()

	router.Run(":" + port)
}
