package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {

	server.GET("/car", getCar)                     //GET, POST, PUT, PATCH, DELETE
	server.GET("/car/:id", getCarById)             // cars/1, cars/2
	server.POST("/car", createCar)                 //POST car details, Payload needs to have a Brand: String, Model: String, Engine: String, Gearbox: String
	server.GET("/car/brand/:brand", getCarByBrand) //GET car by its brand
	server.PUT("/car/:id", updateCar)              //PUT request to update car
	server.DELETE("/car/:id", deleteCar)           //DELETE a car
	server.POST("/car_id", forceCar)               //POST car with ID, Payload needs to have
	server.POST("/signup", signUp)                 //POST to create user with email,password,alias
	server.DELETE("/signup/:id", deleteUser)       //DELETE user with internal ID

}
