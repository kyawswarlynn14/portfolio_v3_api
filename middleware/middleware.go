package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"portfolio/models"
	token "portfolio/tokens"

	"github.com/avct/uasurfer"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Authentication middleware for validating JWT token
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			c.Abort()
			return
		}

		// Parse Bearer token
		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		clientToken := fields[1]
		claims, err := token.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("email", claims.Email)
		c.Set("userId", claims.User_ID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// Authorization middleware to check roles
func Authorization(roles []int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleFromMdw, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Role not found in request context"})
			c.Abort()
			return
		}

		userRole, ok := userRoleFromMdw.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Invalid role format"})
			c.Abort()
			return
		}

		isAuthorized := false
		for _, role := range roles {
			if userRole == role {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func VisitorDataMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.Request.UserAgent()

		deviceType := "Desktop"
		if strings.Contains(strings.ToLower(userAgent), "mobile") {
			deviceType = "Mobile"
		}

		ip := c.ClientIP()
		if ip == "::1" || ip == "127.0.0.1" {
			ip = "127.0.0.1"
		}

		country := "Unknown"
		if net.ParseIP(ip) != nil {
			country, _ = GetCountryFromIP(ip)
		}

		browser, os := ParseUserAgent(userAgent)

		visitorLog := models.VisitorLog{
			ID:        primitive.NewObjectID(),
			Device:    deviceType,
			Country:   country,
			IP:        ip,
			Browser:   browser,
			OS:        os,
			Timestamp: time.Now(),
		}

		c.Set("visitor_log", visitorLog)
		c.Next()
	}
}

// GetCountryFromIP fetches country name from IP
func GetCountryFromIP(ip string) (string, error) {
	url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "Unknown", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Unknown", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "Unknown", err
	}

	if country, ok := result["country_name"].(string); ok {
		return country, nil
	}

	return "Unknown", nil
}

// ParseUserAgent parses the User-Agent string for browser and OS
func ParseUserAgent(userAgent string) (string, string) {
	ua := uasurfer.Parse(userAgent)

	browser := "Unknown"
	if ua.Browser.Name != uasurfer.BrowserUnknown {
		browser = ua.Browser.Name.String()
	}

	os := "Unknown"
	if ua.OS.Name != uasurfer.OSUnknown {
		os = ua.OS.Name.String()
	}

	return browser, os
}
