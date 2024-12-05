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

var CertificateCollection *mongo.Collection = database.CertificateData(database.Client, "Certificates")

func CreateCertificate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var certificate models.Certificate
		if err := c.BindJSON(&certificate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}

		certificate.Certificate_ID = primitive.NewObjectID()
		certificate.Created_At = time.Now()
		certificate.Updated_At = time.Now()

		_, err := CertificateCollection.InsertOne(ctx, certificate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error creating certificate"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Certificate created successfully"})
	}
}

func UpdateCertificate() gin.HandlerFunc {
	return func(c *gin.Context) {
		certificateID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(certificateID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid certificate ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var certificate models.Certificate
		if err := c.BindJSON(&certificate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}

		certificate.Updated_At = time.Now()

		update := bson.M{
			"$set": bson.M{
				"title":      certificate.Title,
				"content":    certificate.Content,
				"image":      certificate.Image,
				"demo_link":  certificate.DemoLink,
				"t1":         certificate.T1,
				"t2":         certificate.T2,
				"updated_at": certificate.Updated_At,
			},
		}

		result, err := CertificateCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error updating certificate", "details": err.Error()})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Certificate not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Certificate updated successfully"})
	}
}

func DeleteCertificate() gin.HandlerFunc {
	return func(c *gin.Context) {
		certificateID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(certificateID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid certificate ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := CertificateCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error deleting certificate"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Certificate not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Certificate deleted successfully"})
	}
}

func GetAllCertificates() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var certificates []models.Certificate
		cursor, err := CertificateCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error retrieving certificates"})
			return
		}

		if err = cursor.All(ctx, &certificates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error decoding certificates"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "certificates": certificates})
	}
}

func GetOneCertificate() gin.HandlerFunc {
	return func(c *gin.Context) {
		certificateID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(certificateID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid certificate ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var certificate models.Certificate
		err = CertificateCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&certificate)
		if err != nil {
			log.Printf("Error retrieving certificate: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error retrieving certificate", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "certificate": certificate})
	}
}
