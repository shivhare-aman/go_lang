package main

import (
	"golang/controllers"
	"golang/dao"
	"golang/initializers"
	"golang/middleware"
	"golang/services"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.InitializeDB()
}

func main() {

	db := initializers.DB
	newUserDao := dao.NewUserDao(db)
	service := services.NewUserService(newUserDao)
	controller := controllers.NewUserController(service)

	authController := controllers.NewAuthController(*newUserDao)

	router := gin.Default()

	router.POST("/users", middleware.RequireAuth("RoleUser", "RoleAdmin"), controller.CreateUser)
	router.GET("/users/:id", middleware.RequireAuth("RoleUser", "RoleAdmin"), controller.GetUserById)
	router.GET("/users", middleware.RequireAuth("RoleUser", "RoleAdmin"), controller.GetAllUsers)
	router.PUT("/users/:id", middleware.RequireAuth("RoleUser", "RoleAdmin"), controller.UpdateUser)
	router.DELETE("/users/:id", middleware.RequireAuth("RoleUser", "RoleAdmin"), controller.DeleteUser)

	router.POST("/signup", controller.CreateUser)
	router.POST("/login", authController.Login)

	router.Run(":8080")
}
