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

var ExpenseCategoryCollection *mongo.Collection = database.PortfolioData(database.Client, "ExpenseCategories")

func CreateExpenseCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDFromMdw, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
			return
		}

		userIDStr, ok := userIDFromMdw.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Invalid user ID format"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var expenseCategory models.ExpenseCategory
		if err := c.BindJSON(&expenseCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}

		expenseCategory.Category_ID = primitive.NewObjectID()
		expenseCategory.User_ID = userIDStr
		expenseCategory.Created_At = time.Now()
		expenseCategory.Updated_At = time.Now()

		_, err := ExpenseCategoryCollection.InsertOne(ctx, expenseCategory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error creating expense category"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Expense Category created successfully"})
	}
}

func UpdateExpenseCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		expenseCategoryID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(expenseCategoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid expense category ID"})
			return
		}

		userIDFromMdw, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
			return
		}

		userIDStr, ok := userIDFromMdw.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Invalid user ID format"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var expenseCategory models.ExpenseCategory
		if err := c.BindJSON(&expenseCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}

		var existingExpenseCategory models.ExpenseItem
		err = ExpenseCategoryCollection.FindOne(ctx, bson.M{"_id": objID, "user_id": userIDStr}).Decode(&existingExpenseCategory)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Expense category is not found or access denied"})
			return
		}

		updateFields := bson.M{}
		if expenseCategory.Title != "" {
			updateFields["title"] = expenseCategory.Title
		}
		if expenseCategory.Description != "" {
			updateFields["description"] = expenseCategory.Description
		}
		if expenseCategory.Type != "" {
			updateFields["type"] = expenseCategory.Type
		}
		if expenseCategory.T1 != "" {
			updateFields["t1"] = expenseCategory.T1
		}
		if expenseCategory.T2 != "" {
			updateFields["t2"] = expenseCategory.T2
		}
		updateFields["updated_at"] = time.Now()

		result, err := ExpenseCategoryCollection.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			bson.M{"$set": updateFields},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error updating expense category", "details": err.Error()})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Expense category not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Expense category updated successfully"})
	}
}

func DeleteExpenseCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		expenseCategoryID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(expenseCategoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid expense category ID"})
			return
		}

		userIDFromMdw, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var expenseCategory models.ExpenseCategory
		err = ExpenseCategoryCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&expenseCategory)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Expense category not found"})
			return
		}

		if expenseCategory.User_ID != userIDFromMdw {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "You are not allowed to delete this expense category"})
			return
		}

		_, err = ExpenseCategoryCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error deleting expense category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Expense category deleted successfully"})
	}
}

func GetAllExpenseCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDFromMdw, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
			return
		}

		userIDStr, ok := userIDFromMdw.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Invalid user ID format"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var expenseCategories []models.ExpenseCategory
		cursor, err := ExpenseCategoryCollection.Find(ctx, bson.M{"user_id": userIDStr})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error retrieving expense categories"})
			return
		}

		if err = cursor.All(ctx, &expenseCategories); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error decoding expense categories"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "categories": expenseCategories})
	}
}

func GetOneExpenseCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		expenseCategoryID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(expenseCategoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid expense category ID"})
			return
		}

		userIDFromMdw, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
			return
		}

		userIDStr, ok := userIDFromMdw.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Invalid user ID format"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var expenseCategory models.ExpenseCategory
		err = ExpenseCategoryCollection.FindOne(ctx, bson.M{"_id": objID, "user_id": userIDStr}).Decode(&expenseCategory)
		if err != nil {
			log.Printf("Error retrieving expense category: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error retrieving expense category", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "category": expenseCategory})
	}
}
