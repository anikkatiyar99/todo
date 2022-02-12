package routes

import (
	controller "github.com/anikkatiyar99/todo/controllers"
	"github.com/anikkatiyar99/todo/middleware"

	"github.com/gin-gonic/gin"
)

//TaskRoutes function
func TaskRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/tasks", controller.GetTasks())
	incomingRoutes.POST("/tasks", controller.CreateTask())
	incomingRoutes.PUT("/tasks/:task_id", controller.UpdateTask())
}
