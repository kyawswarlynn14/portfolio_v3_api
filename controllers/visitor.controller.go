package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"portfolio/database"
	"portfolio/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var VisitorLogsCollection *mongo.Collection = database.PortfolioData(database.Client, "VisitorLogs")

func CreateVisitorLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		visitorLog, exists := c.Get("visitor_log")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Visitor log not found in context"})
			return
		}

		log.Printf("visitor_log >>> %v", visitorLog)

		layoutType := "about_me"
		filter := bson.M{"type": layoutType}

		update := bson.M{
			"$inc": bson.M{"data.view_count": 1},
		}

		_, err := LayoutCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error updating view count in layout"})
			return
		}

		_, err = VisitorLogsCollection.InsertOne(ctx, visitorLog)
		if err != nil {
			log.Printf("error >>> %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error saving visitor log"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Visitor log saved successfully"})
	}
}

func GetAllVisitorLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := VisitorLogsCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error fetching visitor logs"})
			return
		}
		defer cursor.Close(ctx)

		var logs []models.VisitorLog
		if err := cursor.All(ctx, &logs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error decoding visitor logs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": logs})
	}
}

func GetVisitorLogByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid ID format"})
			return
		}

		var log models.VisitorLog
		err = VisitorLogsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&log)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Visitor log not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error fetching visitor log"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": log})
	}
}

func DeleteVisitorLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid ID format"})
			return
		}

		result, err := VisitorLogsCollection.DeleteOne(ctx, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error deleting visitor log"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Visitor log not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Visitor log deleted successfully"})
	}
}
