//@title Login API
//@description API para autenticação de usuários em GOlang
//@version 1.0
//@host localhost:8080
//@BasePath /

package main

import (
	_ "src/docs"
	"src/pkg/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	app := gin.Default()

	app.GET("/getUsers", controller.GetAllUsers)
	app.POST("/login", controller.Login)

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.Run("localhost:8080")

}
