package controllers

import (
	"net/http"
	"os"
	"portfolio/models"
	generate "portfolio/tokens"

	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		EMAIL := os.Getenv("EMAIL")
		PASSWORD := os.Getenv("PASSWORD")

		var loginDetails models.LoginDetails
		if err := c.BindJSON(&loginDetails); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}

		if loginDetails.Email == EMAIL && loginDetails.Password == PASSWORD {
			token, err := generate.TokenGenerator(loginDetails.Email, "", 0)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to generate token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true, "token": token})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Invalid email or password"})
		}
	}
}
