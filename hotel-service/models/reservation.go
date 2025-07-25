package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HotelReservation struct {
    ID        uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
    HotelID   uuid.UUID `gorm:"type:uuid;not null" json:"hotel_id"`
	Hotel     Hotel     `gorm:"foreignKey:HotelID" json:"-"`
    StartDate string `gorm:"type:date;not null" json:"start_date"`
    EndDate   string `gorm:"type:date;not null" json:"end_date"`
    RoomCount int       `gorm:"not null" json:"room_count"`
}

func (r *HotelReservation) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}