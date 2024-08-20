package main

import (
	"gooooo/db"
	"gooooo/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	server.GET("/car", getCar)
	server.GET("/car/:id", getCarById)
	server.POST("/car", createCar)

	server.Run(":8080")
}

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
