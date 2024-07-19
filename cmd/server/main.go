package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Inisialisasi koneksi database
	dsn := "host=localhost user=postgres password=postgres dbname=task_management port=5432 sslmode=disable"
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	// Inisialisasi router Gin
	r := gin.Default()

	// Definisikan rute dasar
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Task Management System",
		})
	})

	// Jalankan server
	r.Run(":8080")
}
