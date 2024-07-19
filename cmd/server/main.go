package main

import (
	"log"
	"net/http"

	"github.com/andrMaulana/Go-Task-Management-System/internal/database"
	"github.com/andrMaulana/Go-Task-Management-System/internal/handler"
	"github.com/andrMaulana/Go-Task-Management-System/internal/middleware"
	"github.com/andrMaulana/Go-Task-Management-System/internal/models"
	"github.com/andrMaulana/Go-Task-Management-System/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Inisialisasi koneksi database
	dsn := "host=localhost user=postgres password=postgres dbname=task_management port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}
	// Jalankan migrasi
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations")
	}

	// Inisialisasi service
	tokenService := service.NewTokenService()
	userService := service.NewUserService(db)
	projectService := service.NewProjectService(db)
	taskService := service.NewTaskService(db)

	// Inisialisasi handler
	userHandler := handler.NewUserHandler(userService, tokenService)
	projectHandler := handler.NewProjectHandler(projectService)
	taskHandler := handler.NewTaskHandler(taskService)

	// Inisialisasi router Gin
	r := gin.Default()

	// Definisikan rute dasar
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Task Management System",
		})
	})

	// Definisikan rute untuk user
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/logout", middleware.AuthMiddleware(tokenService), userHandler.Logout)

	// Grup rute yang memerlukan autentikasi
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware(tokenService))
	{
		// Tambahkan rute yang memerlukan autentikasi di sini
		authorized.GET("/protected", func(c *gin.Context) {
			userId := c.MustGet("user_id").(float64)
			c.JSON(http.StatusOK, models.Response{
				Meta: models.Meta{
					Code:    http.StatusOK,
					Status:  "OK",
					Message: "This is a protected route",
				},
				Data: gin.H{"user_id": userId},
			})
		})
	}

	// Jalankan server
	r.Run(":8080")
}
