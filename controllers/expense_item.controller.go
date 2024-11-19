package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"portfolio/database"
	"portfolio/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ExpenseItemCollection *mongo.Collection = database.ExpenseCategoryData(database.Client, "ExpenseItems")

func CreateExpenseItem() gin.HandlerFunc {
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

		var expenseItem models.ExpenseItem
		if err := c.BindJSON(&expenseItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}

		expenseItem.Item_ID = primitive.NewObjectID()
		expenseItem.User_ID = userIDStr
		expenseItem.Created_At = time.Now()
		expenseItem.Updated_At = time.Now()

		_, err := ExpenseItemCollection.InsertOne(ctx, expenseItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error creating expense item"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Expense item created successfully"})
	}
}

func UpdateExpenseItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		itemID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(itemID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid expense item ID"})
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

		var expenseItem models.ExpenseItem
		if err := c.BindJSON(&expenseItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}

		var existingItem models.ExpenseItem
		err = ExpenseItemCollection.FindOne(ctx, bson.M{"_id": objID, "user_id": userIDStr}).Decode(&existingItem)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Expense item not found or access denied"})
			return
		}

		updateFields := bson.M{
			"updated_at": time.Now(),
		}
		if expenseItem.Title != "" {
			updateFields["title"] = expenseItem.Title
		}
		if expenseItem.Remark != "" {
			updateFields["remark"] = expenseItem.Remark
		}
		if expenseItem.Amount != 0 {
			updateFields["amount"] = expenseItem.Amount
		}
		if expenseItem.T1 != "" {
			updateFields["t1"] = expenseItem.T1
		}
		if expenseItem.T2 != "" {
			updateFields["t2"] = expenseItem.T2
		}

		result, err := ExpenseItemCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateFields})
		if err != nil || result.MatchedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error updating expense item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Expense item updated successfully"})
	}
}

func DeleteExpenseItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		itemID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(itemID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid expense item ID"})
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

		result, err := ExpenseItemCollection.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userIDStr})
		if err != nil || result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Expense item not found or access denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Expense item deleted successfully"})
	}
}

func GetAllIncomes() gin.HandlerFunc {
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

		yearQuery := c.Query("year")
		if yearQuery == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Year is required"})
			return
		}

		year, err := strconv.Atoi(yearQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid year format"})
			return
		}

		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(1, 0, 0).Add(-time.Second)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var incomes []models.ExpenseItem
		cursor, err := ExpenseItemCollection.Find(ctx, bson.M{
			"user_id": userIDStr,
			"type":    "001",
			"created_at": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error retrieving incomes"})
			return
		}

		if err = cursor.All(ctx, &incomes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error decoding incomes"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "incomes": incomes})
	}
}

func GetAllOutcomes() gin.HandlerFunc {
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

		yearQuery := c.Query("year")
		if yearQuery == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Year is required"})
			return
		}

		year, err := strconv.Atoi(yearQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid year format"})
			return
		}

		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(1, 0, 0).Add(-time.Second)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var outcomes []models.ExpenseItem
		cursor, err := ExpenseItemCollection.Find(ctx, bson.M{
			"user_id": userIDStr,
			"type":    "002",
			"created_at": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error retrieving outcomes"})
			return
		}

		if err = cursor.All(ctx, &outcomes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Error decoding outcomes"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "outcomes": outcomes})
	}
}
