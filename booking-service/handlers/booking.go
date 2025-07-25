package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"booking-service/config"
	"booking-service/kafka"
	"booking-service/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetDB(db *gorm.DB) {
	DB = db
}

type createBookingRequest struct {
	HotelID   uuid.UUID `json:"hotel_id" binding:"required"`
	FlightID  uuid.UUID `json:"flight_id" binding:"required"`
	RoomCount int       `json:"room_count" binding:"required"`
	SeatCount int       `json:"seat_count" binding:"required"`
	StartDate string    `json:"start_date" binding:"required"`
	EndDate   string    `json:"end_date" binding:"required"`
}

func CreateBooking(c *gin.Context) {
	cfg := config.LoadConfig()
	var req createBookingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking := models.Booking{
		HotelID:   req.HotelID,
		FlightID:  req.FlightID,
		RoomCount: req.RoomCount,
		SeatCount: req.SeatCount,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    "PENDING",
	}

	if result := DB.Create(&booking); result.Error != nil {
		log.Printf("Error creating booking: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	bookingID := booking.ID.String()

	event := kafka.BookingRequestedEvent{
		BookingID: bookingID,
		HotelID:   req.HotelID,
		FlightID:  req.FlightID,
		RoomCount: req.RoomCount,
		SeatCount: req.SeatCount,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
	eventBytes, _ := json.Marshal(event)

	err := kafka.ProduceMessage(cfg.KafkaRequestTopic, []byte(bookingID), eventBytes)
	if err != nil {
		log.Printf("Error publishing booking request to Kafka: %v", err)
		booking.Status = "FAILED"
		DB.Save(&booking)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate booking process"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":    "Booking request received, processing initiated",
		"booking_id": bookingID,
	})

	c.JSON(http.StatusCreated, booking)
}
