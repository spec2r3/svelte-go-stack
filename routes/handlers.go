package routes

import (
	"gooooo/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCar(context *gin.Context) {
	cars, err := models.GetAllCars()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve cars"})
		return
	}
	context.JSON(http.StatusOK, cars)
}

func getCarById(context *gin.Context) {
	carId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid car ID"})
		return
	}

	car, err := models.GetCarByID(carId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve the car"}) //Check status later
		return
	}
	if car == nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Car not found"})
		return
	}

	context.JSON(http.StatusOK, car)
}

func createCar(context *gin.Context) {
	var car models.Car
	err := context.ShouldBindJSON(&car)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid format"})
		return
	}

	err = car.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save car"})
		return
	}

	context.JSON(http.StatusCreated, car)
}

func getCarByBrand(c *gin.Context) {
	brand := c.Param("brand")

	// Fetch car details by brand
	cars, err := models.GetCarsByBrand(brand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve cars"})
		return
	}
	if len(cars) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No cars found with the given brand"})
		return
	}

	c.JSON(http.StatusOK, cars)
}

func updateCar(context *gin.Context) {
	carId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse car id."})
		return
	}

	_, err = models.GetCarByID(carId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the car."})
		return
	}

	var updatedCar models.Car
	err = context.ShouldBindJSON(&updatedCar)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid format"})
		return
	}

	updatedCar.ID = carId
	err = updatedCar.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update car"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Car updated successfully"})

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
