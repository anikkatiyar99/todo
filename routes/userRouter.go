package routes

import (
	controller "github.com/anikkatiyar99/todo/controllers"
	"github.com/anikkatiyar99/todo/middleware"

	"github.com/gin-gonic/gin"
)

//UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}
