package routes

import (
	"gooooo/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCar(c *gin.Context) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	if page < 1 && pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters"})
		return
	}

	cars, carCount, err := models.GetAllCars(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve cars"})
		return
	}

	totalPages := (carCount + pageSize - 1) / pageSize

	response := gin.H{
		"cars":        cars,
		"totalCount":  carCount,
		"totalPages":  totalPages,
		"currentPage": page,
		"pageSize":    pageSize,
	}

	c.JSON(http.StatusOK, response)
}

func getCarById(c *gin.Context) {
	carId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid car ID"})
		return
	}

	car, err := models.GetCarByID(carId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve the car"})
		return
	}
	if car == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Car not found"})
		return
	}

	c.JSON(http.StatusOK, car)
}

func createCar(c *gin.Context) {
	var cars []models.Car
	err := c.ShouldBindJSON(&cars)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid format"})
		return
	}

	var createdModels []string
	var duplicateModels []string

	for _, car := range cars {
		exists, err := models.CarCountChecker(car.Brand, car.Model)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error checking for duplicates"})
			return
		}
		if exists {
			duplicateModels = append(duplicateModels, car.Model)
			continue
		}

		err = car.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save car"})
			return
		}
		createdModels = append(createdModels, car.Model)
	}

	response := gin.H{
		"created":    createdModels,
		"duplicates": duplicateModels,
	}

	c.JSON(http.StatusCreated, response)
}

func getCarByBrand(c *gin.Context) {
	brand := c.Param("brand")

	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	page := 1
	pageSize := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
		pageSize = ps
	}

	if page < 1 && pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters"})
		return
	}

	cars, totalCount, err := models.GetCarsByBrand(brand, page, pageSize)
	if err != nil {
		if len(cars) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No cars found with the given brand"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve cars"})
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"cars":        cars,
		"totalCount":  totalCount,
		"totalPages":  totalPages,
		"currentPage": page,
	})
}

func updateCar(c *gin.Context) {
	carId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse car id."})
		return
	}

	_, err = models.GetCarByID(carId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the car."})
		return
	}

	var updatedCar models.Car
	err = c.ShouldBindJSON(&updatedCar)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid format"})
		return
	}

	updatedCar.ID = carId
	err = updatedCar.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update car"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Car updated successfully"})

}

func deleteCar(c *gin.Context) {

	carId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid car ID"})
		return
	}

	err = models.DeleteCarById(carId)
	if err != nil {
		if err.Error() == "no car found with ID" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Car not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the car"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car deleted successfully"})
}

func forceCar(c *gin.Context) {

	var car models.Car
	err := c.ShouldBindJSON(&car)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid format"})
		return
	}

	if car.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID is required"})
		return
	}

	err = car.Force()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save car"})
		return
	}

	c.JSON(http.StatusCreated, car)

}
