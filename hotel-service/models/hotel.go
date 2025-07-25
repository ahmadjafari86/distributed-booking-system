package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hotel struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	TotalRooms int       `gorm:"not null" json:"total_rooms"`
	Reservations []HotelReservation `gorm:"foreignKey:HotelID"`
}

func (h *Hotel) BeforeCreate(tx *gorm.DB) (err error) {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return
}