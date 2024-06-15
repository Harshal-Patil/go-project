package employee

import (
	"fmt"
	"log"
	"net/http"

	"go-project/controllers"

	"github.com/gin-gonic/gin"
)

func DeleteData(c *gin.Context) {
	db := controllers.GetDB()
	redisClient := controllers.GetRedisClient()
	ctx := controllers.GetContext()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing employee ID"})
		return
	}

	// Validate the ID

	result, err := db.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting data from database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data from database"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	redisKey := fmt.Sprintf("employee:%s", id)
	if err := redisClient.Del(ctx, redisKey).Err(); err != nil {
		log.Printf("Error deleting data from Redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data from Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Record deleted successfully"})
}
