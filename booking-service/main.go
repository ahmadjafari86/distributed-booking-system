package main

import (
	"fmt"
	"log"
	"net/http"

	"booking-service/config"
	"booking-service/handlers"
	"booking-service/kafka"
	"booking-service/models"
	"booking-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.Booking{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database schema: %v", err)
	}
	log.Println("Database migration complete.")

	handlers.SetDB(db)
	services.SetDB(db)

	kafka.InitProducer(cfg.KafkaBrokerAddress)

	kafka.StartConsumer(cfg.KafkaBrokerAddress, cfg.KafkaGroupID, cfg.KafkaResultTopic, services.HandleBookingResults)

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/booking/book-trip", handlers.CreateBooking)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "booking service is running!"})
	})

	serverAddress := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on port %s...", cfg.ServerPort)
	log.Fatal(router.Run(serverAddress))
}
