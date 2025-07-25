package handlers

import (
	"fmt"
	"net/http"
	"time"

	"hotel-service/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


func CreateHotelReservation(c *gin.Context) {
	var reqReservation models.HotelReservation
	var reservation models.HotelReservation

	if err := c.ShouldBindJSON(&reqReservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const dateFormat = "2006-01-02" // YYYY-MM-DD format

	parsedStartDate, err := time.Parse(dateFormat, reqReservation.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Please use YYYY-MM-DD."})
		return
	}

	parsedEndDate, err := time.Parse(dateFormat, reqReservation.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Please use YYYY-MM-DD."})
		return
	}

	reservation.HotelID = reqReservation.HotelID
	reservation.StartDate = reqReservation.StartDate
	reservation.EndDate = reqReservation.EndDate
	reservation.RoomCount = reqReservation.RoomCount

	if reservation.RoomCount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room count must be greater than 0."})
		return
	}

	if parsedStartDate.After(parsedEndDate) || parsedStartDate.Equal(parsedEndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date."})
		return
	}

	var hotel models.Hotel
	if result := DB.First(&hotel, "id = ?", reservation.HotelID); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotel details", "details": result.Error.Error()})
		}
		return
	}

	for d := parsedStartDate; d.Before(parsedEndDate); d = d.AddDate(0, 0, 1) {
		var reservedRooms int
		var overlappingReservations []models.HotelReservation

		if result := DB.Where("hotel_id = ? AND start_date <= ? AND end_date > ?", hotel.ID, d, d).Find(&overlappingReservations); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to check availability for %s", d.Format(dateFormat)), "details": result.Error.Error()})
			return
		}

		for _, res := range overlappingReservations {
			reservedRooms += res.RoomCount
		}

		if (reservedRooms + reservation.RoomCount) > hotel.TotalRooms {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Rooms not available for %s. Only %d rooms free.", d.Format(dateFormat), hotel.TotalRooms-reservedRooms)})
			return
		}
	}

	if reservation.ID == uuid.Nil {
		reservation.ID = uuid.New()
	}

	if result := DB.Create(&reservation); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

func ListHotelReservations(c *gin.Context) {
	var reservations []models.HotelReservation
	if result := DB.Preload("Hotel").Find(&reservations); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reservations", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

func CancelHotelReservation(c *gin.Context) {
	reservationIDStr := c.Param("id")
	reservationID, err := uuid.Parse(reservationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID format"})
		return
	}

	var reservation models.HotelReservation
	if result := DB.First(&reservation, "id = ?", reservationID); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reservation", "details": result.Error.Error()})
		}
		return
	}

	if result := DB.Delete(&reservation); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel reservation", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation cancelled successfully"})
}