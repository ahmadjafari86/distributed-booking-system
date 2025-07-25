package main

import (
	"fmt"
	"log"
	"net/http"

	"hotel-service/config"
	"hotel-service/handlers"
	"hotel-service/kafka"
	"hotel-service/models"
	"hotel-service/services"

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

	err = db.AutoMigrate(&models.Hotel{}, &models.HotelReservation{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database schema: %v", err)
	}
	log.Println("Database migration complete.")

	handlers.SetDB(db)

	router := gin.Default()

	kafka.InitProducer(cfg.KafkaBrokerAddress)

	kafka.StartConsumer(cfg.KafkaBrokerAddress, cfg.KafkaGroupID, cfg.KafkaRequestTopic, services.HandleBookingResults)

	api := router.Group("/api/v1")
	{
		api.POST("/hotels", handlers.CreateHotel)
		api.GET("/hotels", handlers.GetHotels)
		api.GET("/hotels/free-rooms", handlers.GetFreeRoomsByDate)
		api.GET("/hotels/:hotel_id/reservations", handlers.ListReservationsByHotelID)
		api.POST("/reservations", handlers.CreateHotelReservation)
		api.GET("/reservations", handlers.ListHotelReservations)
		api.DELETE("/reservations/:id", handlers.CancelHotelReservation)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hotel reservation service is running!"})
	})

	serverAddress := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on port %s...", cfg.ServerPort)
	log.Fatal(router.Run(serverAddress))
}
