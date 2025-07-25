package kafka

import "github.com/google/uuid"

type BookingRequest struct {
	HotelID   string `json:"hotel_id"`
	FlightID  string `json:"flight_id"`
	RoomCount int    `json:"room_count"`
	SeatCount int    `json:"seat_count"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type BookingRequestedEvent struct {
	BookingID string    `json:"booking_id"`
	HotelID   uuid.UUID `json:"hotel_reservation_id"`
	FlightID  uuid.UUID `json:"flight_booking_id"`
	RoomCount int       `json:"room_count"`
	SeatCount int       `json:"seat_count"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
}

type BookingResultEvent struct {
	BookingID   string `json:"booking_id"`
	ServiceType string `json:"service_type"`
	ServiceID   string `json:"service_id"`
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
}

type CancelBookingEvent struct {
	BookingID string `json:"booking_id"`
	ServiceID string `json:"service_id"`
}
