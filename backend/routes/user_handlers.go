package routes

import (
	"gooooo/models"
	"gooooo/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func signUp(context *gin.Context) {

	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid format"})
		return
	}

	apiKey, err := utils.GenerateKey()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate API key"})
		return
	}

	user.APIKey = apiKey

	err = user.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create User"})
		return
	}

	context.JSON(http.StatusCreated, "User created successfully")
}

func deleteUser(c *gin.Context) {

	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	err = models.DeleteUserById(userId)
	if err != nil {
		if err.Error() == "no user found with ID" {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the User"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// func signIn(c *gin.Context) {

// 	var credentials models.SignIn

// 	err := c.ShouldBindJSON(&car)

// 	userName, err := c.Name
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid car ID"})
// 		return
// 	}

// 	err = models.DeleteUserById(userId)
// 	if err != nil {
// 		if err.Error() == "no user found with ID" {
// 			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the User"})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
// }

func getUser(c *gin.Context) {
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

	users, userCount, err := models.GetAllUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve cars"})
		return
	}

	var userResp []models.UserResponse

	for _, user := range users {
		userResp = append(userResp, models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Alias: user.Alias,
			Admin: user.Admin,
		})
	}

	totalPages := (userCount + pageSize - 1) / pageSize

	response := gin.H{
		"users":       userResp,
		"totalCount":  userCount,
		"totalPages":  totalPages,
		"currentPage": page,
		"pageSize":    pageSize,
	}

	c.JSON(http.StatusOK, response)
}

func getUserById(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve the user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	response := gin.H{
		"ID":    user.ID,
		"Email": user.Email,
		"Alias": user.Alias,
	}

	c.JSON(http.StatusOK, response)
}

func updateUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user id."})
		return
	}

	_, err = models.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the car."})
		return
	}

	var updatedUser models.Car
	err = c.ShouldBindJSON(&updatedUser)

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
