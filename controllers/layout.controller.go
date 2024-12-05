package controllers

import (
	"context"
	"net/http"
	"time"

	"portfolio/database"
	"portfolio/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var LayoutCollection *mongo.Collection = database.LayoutData(database.Client, "Layouts")

func ManageLayout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		layoutType := c.Query("type")

		switch c.Request.Method {
		case http.MethodGet:
			getLayout(ctx, c, layoutType)
		case http.MethodPost:
			createLayout(ctx, c, layoutType)
		case http.MethodPut:
			updateLayout(ctx, c, layoutType)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid method"})
		}
	}
}

func createLayout(ctx context.Context, c *gin.Context, layoutType string) {
	var layoutData interface{}
	var err error

	switch layoutType {
	case "about_me":
		var aboutMe models.AboutMe
		err = c.BindJSON(&aboutMe)
		layoutData = aboutMe
	case "service_info":
		var serviceInfo models.ServiceInfo
		err = c.BindJSON(&serviceInfo)
		layoutData = serviceInfo
	case "project_info":
		var projectInfo models.ProjectInfo
		err = c.BindJSON(&projectInfo)
		layoutData = projectInfo
	case "blog_info":
		var blogInfo models.Blog
		err = c.BindJSON(&blogInfo)
		layoutData = blogInfo
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid layout type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	_, err = LayoutCollection.InsertOne(ctx, bson.M{"type": layoutType, "data": layoutData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error creating layout"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": layoutType + " created successfully"})
}

func updateLayout(ctx context.Context, c *gin.Context, layoutType string) {
	var layoutData interface{}
	var err error

	switch layoutType {
	case "about_me":
		var aboutMe models.AboutMe
		err = c.BindJSON(&aboutMe)
		layoutData = aboutMe
	case "service_info":
		var serviceInfo models.ServiceInfo
		err = c.BindJSON(&serviceInfo)
		layoutData = serviceInfo
	case "project_info":
		var projectInfo models.ProjectInfo
		err = c.BindJSON(&projectInfo)
		layoutData = projectInfo
	case "blog_info":
		var blogInfo models.Blog
		err = c.BindJSON(&blogInfo)
		layoutData = blogInfo
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid layout type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	update := bson.M{"$set": bson.M{"data": layoutData}}

	result, err := LayoutCollection.UpdateOne(ctx, bson.M{"type": layoutType}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error updating layout"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": layoutType + " not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": layoutType + " updated successfully"})
}

func getLayout(ctx context.Context, c *gin.Context, layoutType string) {
	type LayoutData struct {
		ID   primitive.ObjectID `json:"_id" bson:"_id"`
		Type string             `json:"type" bson:"type"`
		Data interface{}        `json:"data" bson:"data"`
	}

	var layoutData LayoutData
	layoutData.Type = layoutType

	switch layoutType {
	case "about_me":
		var aboutMe models.AboutMe
		err := LayoutCollection.FindOne(ctx, bson.M{"type": layoutType}).Decode(&struct {
			ID   primitive.ObjectID `json:"_id" bson:"_id"`
			Type string             `json:"type" bson:"type"`
			Data *models.AboutMe    `json:"data" bson:"data"`
		}{
			ID:   layoutData.ID,
			Type: layoutData.Type,
			Data: &aboutMe,
		})
		layoutData.Data = aboutMe
		if err != nil {
			handleError(c, err, layoutType)
			return
		}
	case "service_info":
		var serviceInfo models.ServiceInfo
		err := LayoutCollection.FindOne(ctx, bson.M{"type": layoutType}).Decode(&struct {
			ID   primitive.ObjectID   `json:"_id" bson:"_id"`
			Type string               `json:"type" bson:"type"`
			Data *models.ServiceInfo  `json:"data" bson:"data"`
		}{
			ID:   layoutData.ID,
			Type: layoutData.Type,
			Data: &serviceInfo,
		})
		layoutData.Data = serviceInfo
		if err != nil {
			handleError(c, err, layoutType)
			return
		}
	case "project_info":
		var projectInfo models.ProjectInfo
		err := LayoutCollection.FindOne(ctx, bson.M{"type": layoutType}).Decode(&struct {
			ID   primitive.ObjectID   `json:"_id" bson:"_id"`
			Type string               `json:"type" bson:"type"`
			Data *models.ProjectInfo  `json:"data" bson:"data"`
		}{
			ID:   layoutData.ID,
			Type: layoutData.Type,
			Data: &projectInfo,
		})
		layoutData.Data = projectInfo
		if err != nil {
			handleError(c, err, layoutType)
			return
		}
	case "blog_info":
		var blogInfo models.Blog
		err := LayoutCollection.FindOne(ctx, bson.M{"type": layoutType}).Decode(&struct {
			ID   primitive.ObjectID `json:"_id" bson:"_id"`
			Type string             `json:"type" bson:"type"`
			Data *models.Blog       `json:"data" bson:"data"`
		}{
			ID:   layoutData.ID,
			Type: layoutData.Type,
			Data: &blogInfo,
		})
		layoutData.Data = blogInfo
		if err != nil {
			handleError(c, err, layoutType)
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid layout type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "layout": layoutData})
}

func handleError(c *gin.Context, err error, layoutType string) {
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": layoutType + " not found"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error retrieving " + layoutType})
	}
}

