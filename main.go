package main

import (
	"github.com/gin-gonic/gin"
	"src/pkg/controller"
)

func main() {
	app := gin.Default()

	app.GET("/getUsers", controller.GetAllUsers)
	app.POST("/login", controller.Login)
	app.Run("localhost:8080")

}