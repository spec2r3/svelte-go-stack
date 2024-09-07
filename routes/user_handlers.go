package routes

import (
	"gooooo/models"
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

	// Generate the API key
	apiKey, err := models.GenerateKey()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate API key"})
		return
	}

	// Assign the generated API key to the user struct
	user.APIKey = apiKey

	err = user.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create User"})
		return
	}

	context.JSON(http.StatusCreated, user)
}

func deleteUser(c *gin.Context) {

	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid car ID"})
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
