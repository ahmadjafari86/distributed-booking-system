package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending      BookingStatus = "PENDING"
	BookingStatusHotelBooked  BookingStatus = "HOTEL_BOOKED"
	BookingStatusFlightBooked BookingStatus = "FLIGHT_BOOKED"
	BookingStatusSuccess      BookingStatus = "SUCCESS"
	BookingStatusFailed       BookingStatus = "FAILED"
	BookingStatusCancelled    BookingStatus = "CANCELLED"
)

type Booking struct {
	ID                      uuid.UUID     `gorm:"primaryKey;type:uuid" json:"id"`
	HotelID                 uuid.UUID     `gorm:"type:uuid; not null" json:"hotel_reservation_id"`
	FlightID                uuid.UUID     `gorm:"type:uuid; not null" json:"flight_booking_id"`
	RoomCount               int           `gorm:"type:int; not null" json:"room_count"`
	SeatCount               int           `gorm:"type:int; not null" json:"seat_count"`
	StartDate               string        `gorm:"type:date;not null" json:"start_date"`
	EndDate                 string        `gorm:"type:date;not null" json:"end_date"`
	Status                  BookingStatus `gorm:"type:varchar(50);not null" json:"status"`
	HotelReservationDetails *string       `gorm:"type:jsonb" json:"hotel_reservation_details"`
	FlightBookingDetails    *string       `gorm:"type:jsonb" json:"flight_booking_details"`
	CreatedAt               time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}
