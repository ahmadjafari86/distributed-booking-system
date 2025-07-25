package services

import (
	kafkaBroker "booking-service/kafka"
	"booking-service/models"
	"encoding/json"
	"log"

	kafka "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetDB(db *gorm.DB) {
	DB = db
}

func HandleBookingResults(msg kafka.Message) {
	var event kafkaBroker.BookingResultEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Error unmarshalling booking result event: %v", err)
		return
	}

	var booking models.Booking
	if err := DB.Where("booking_id = ?", event.BookingID).First(&booking).Error; err != nil {
		log.Printf("Booking with ID %s not found: %v", event.BookingID, err)
		return
	}

	if err := DB.Save(&booking).Error; err != nil {
		log.Printf("Failed to update booking %s status: %v", event.BookingID, err)
	} else {
		log.Printf("Booking status updated to: %s", booking.Status)
	}
}
