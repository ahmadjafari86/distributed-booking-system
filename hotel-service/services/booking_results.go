package services

import (
	"encoding/json"
	kafkaBroker "hotel-service/kafka"
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
	log.Printf("Received booking result event: %+v", event)
}
