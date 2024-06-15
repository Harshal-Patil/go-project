package employee

import (
	"fmt"
	"log"
	"net/http"

	"go-project/controllers"

	"github.com/gin-gonic/gin"
)

func EditData(c *gin.Context) {
	db := controllers.GetDB()
	redisClient := controllers.GetRedisClient()
	ctx := controllers.GetContext()

	id := c.Param("id")
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		CompanyName string `json:"company_name"`
		Address     string `json:"address"`
		City        string `json:"city"`
		County      string `json:"county"`
		Postal      string `json:"postal"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Web         string `json:"web"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result, err := db.Exec("UPDATE employees SET first_name = ?, last_name = ?, company_name = ?, address = ?, city = ?, county = ?, postal = ?, phone = ?, email = ?, web = ? WHERE id = ?",
		input.FirstName, input.LastName, input.CompanyName, input.Address, input.City, input.County, input.Postal, input.Phone, input.Email, input.Web, id)
	if err != nil {
		log.Printf("Error updating data in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating data in database: %v", err)})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	redisKey := fmt.Sprintf("employee:%s", id)
	if err := redisClient.Del(ctx, redisKey).Err(); err != nil {
		log.Printf("Error deleting cache from Redis: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"status": "Record updated successfully"})
}
