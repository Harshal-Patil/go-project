package employee

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-project/controllers"

	"github.com/gin-gonic/gin"
)

func ViewData(c *gin.Context) {
	db := controllers.GetDB()
	redisClient := controllers.GetRedisClient()
	ctx := controllers.GetContext()

	id := c.Param("id")
	redisKey := fmt.Sprintf("employee:%s", id)
	employee, err := redisClient.HGetAll(ctx, redisKey).Result()
	if err != nil {
		log.Printf("Error retrieving data from Redis: %v", err)
	}
	if len(employee) > 0 {
		c.JSON(http.StatusOK, employee)
		return
	}

	var (
		firstName, lastName, companyName, address, city, county, postal, phone, email, web string
	)
	err = db.QueryRow("SELECT first_name, last_name, company_name, address, city, county, postal, phone, email, web FROM employees WHERE id = ?", id).Scan(&firstName, &lastName, &companyName, &address, &city, &county, &postal, &phone, &email, &web)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		} else {
			log.Printf("Error retrieving data from database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving data from database"})
		}
		return
	}
	employee = map[string]string{
		"id":           id,
		"first_name":   firstName,
		"last_name":    lastName,
		"company_name": companyName,
		"address":      address,
		"city":         city,
		"county":       county,
		"postal":       postal,
		"phone":        phone,
		"email":        email,
		"web":          web,
	}
	redisClient.HMSet(ctx, redisKey, employee)
	redisClient.Expire(ctx, redisKey, 5*time.Minute)
	c.JSON(http.StatusOK, employee)
}
