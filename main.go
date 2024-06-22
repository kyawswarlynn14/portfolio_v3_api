package main

import (
	"log"
	"os"
	"portfolio/middleware"
	"portfolio/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	publicRoutes := router.Group("/")

	authenticatedRoutes := router.Group("/")
	authenticatedRoutes.Use(middleware.Authentication())

	routes.AuthRoutes(publicRoutes, authenticatedRoutes)
	routes.LayoutRoutes(publicRoutes, authenticatedRoutes)
	routes.CertificateRoutes(publicRoutes, authenticatedRoutes)
	routes.ServiceRoutes(publicRoutes, authenticatedRoutes)
	routes.ProjectRoutes(publicRoutes, authenticatedRoutes)
	routes.EmailRoutes(publicRoutes, authenticatedRoutes)

	log.Fatal(router.Run(":" + port))
}
