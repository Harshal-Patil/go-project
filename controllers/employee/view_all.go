package employee

import (
	"fmt"
	"log"
	"net/http"

	"go-project/controllers"

	"github.com/gin-gonic/gin"
)

func ViewAllData(c *gin.Context) {
	db := controllers.GetDB()

	rows, err := db.Query("SELECT id, first_name, last_name, company_name, address, city, county, postal, phone, email, web FROM employees")
	if err != nil {
		log.Printf("Error retrieving data from database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error retrieving data from database: %v", err)})
		return
	}
	defer rows.Close()

	var employees []map[string]string
	for rows.Next() {
		var (
			id, firstName, lastName, companyName, address, city, county, postal, phone, email, web string
		)
		err := rows.Scan(&id, &firstName, &lastName, &companyName, &address, &city, &county, &postal, &phone, &email, &web)
		if err != nil {
			log.Printf("Error scanning data from database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error scanning data from database: %v", err)})
			return
		}
		employee := map[string]string{
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
		employees = append(employees, employee)
	}

	c.JSON(http.StatusOK, employees)
}
