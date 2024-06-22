package main

import (
	"log"
	"os"
	"portfolio/middleware"
	"portfolio/routes"

	"github.com/gin-contrib/cors"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://kyawswarlynn.vercel.app",
			"https://kyawswarlynn.netlify.app",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	publicRoutes := router.Group("/portfolio/")
	authenticatedRoutes := router.Group("/portfolio/")
	authenticatedRoutes.Use(middleware.Authentication())

	routes.AuthRoutes(publicRoutes, authenticatedRoutes)
	routes.LayoutRoutes(publicRoutes, authenticatedRoutes)
	routes.CertificateRoutes(publicRoutes, authenticatedRoutes)
	routes.ServiceRoutes(publicRoutes, authenticatedRoutes)
	routes.ProjectRoutes(publicRoutes, authenticatedRoutes)
	routes.EmailRoutes(publicRoutes, authenticatedRoutes)

	log.Fatal(router.Run(":" + port))
}
