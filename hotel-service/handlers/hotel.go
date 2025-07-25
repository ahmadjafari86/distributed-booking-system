package handlers

import (
	"net/http"
	"time"

	"hotel-service/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetDB(db *gorm.DB) {
	DB = db
}



func CreateHotel(c *gin.Context) {
	var hotel models.Hotel

	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if hotel.ID == uuid.Nil {
		hotel.ID = uuid.New()
	}

	if result := DB.Create(&hotel); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hotel", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, hotel)
}

func GetHotels(c *gin.Context) {
	var hotels []models.Hotel
	if result := DB.Preload("Reservations").Find(&hotels); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hotels", "details": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, hotels)
}


type FreeRoomsResponse struct {
	HotelID    string `json:"hotel_id"`
	HotelName  string `json:"hotel_name"`
	TotalRooms int    `json:"total_rooms"`
	FreeRooms  int    `json:"free_rooms"`
}

func GetFreeRoomsByDate(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter is required (YYYY-MM-DD)"})
		return
	}

	queryDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Please use YYYY-MM-DD."})
		return
	}

	var hotels []models.Hotel
	if result := DB.Find(&hotels); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotels", "details": result.Error.Error()})
		return
	}

	var response []FreeRoomsResponse
	for _, hotel := range hotels {
		var reservedRooms int
		var reservations []models.HotelReservation

		if result := DB.Where("hotel_id = ? AND start_date <= ? AND end_date > ?", hotel.ID, queryDate, queryDate).Find(&reservations); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations for hotel " + hotel.Name, "details": result.Error.Error()})
			return
		}

		for _, res := range reservations {
			reservedRooms += res.RoomCount
		}

		freeRooms := hotel.TotalRooms - reservedRooms
		if freeRooms < 0 {
			freeRooms = 0
		}

		response = append(response, FreeRoomsResponse{
			HotelID:    hotel.ID.String(),
			HotelName:  hotel.Name,
			TotalRooms: hotel.TotalRooms,
			FreeRooms:  freeRooms,
		})
	}

	c.JSON(http.StatusOK, response)
}

func ListReservationsByHotelID(c *gin.Context) {
	hotelIDStr := c.Param("hotel_id")
	hotelID, err := uuid.Parse(hotelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	var hotel models.Hotel
	if result := DB.First(&hotel, "id = ?", hotelID); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hotel details", "details": result.Error.Error()})
		}
		return
	}

	var reservations []models.HotelReservation
	if result := DB.Preload("Hotel").Where("hotel_id = ?", hotelID).Find(&reservations); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reservations for hotel", "details": result.Error.Error()})
		return
	}

	if len(reservations) == 0 {
		c.JSON(http.StatusOK, []models.HotelReservation{})
		return
	}

	c.JSON(http.StatusOK, reservations)
}