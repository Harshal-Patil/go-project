package employee

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"golang.org/x/sync/errgroup"

	"go-project/controllers"
)

func validateExcelFile(file *xlsx.File) error {
	requiredHeaders := []string{"first_name", "last_name", "company_name", "address", "city", "county", "postal", "phone", "email", "web"}

	if len(file.Sheets) == 0 {
		return fmt.Errorf("no sheets found in the Excel file")
	}

	sheet := file.Sheets[0]
	if len(sheet.Rows) == 0 {
		return fmt.Errorf("no rows found in the Excel sheet")
	}

	headerRow := sheet.Rows[0]
	if len(headerRow.Cells) < len(requiredHeaders) {
		return fmt.Errorf("incorrect number of columns in the header row")
	}

	for i, header := range requiredHeaders {
		if headerRow.Cells[i].String() != header {
			return fmt.Errorf("invalid header: expected %s but got %s", header, headerRow.Cells[i].String())
		}
	}

	return nil
}

func UploadExcel(c *gin.Context) {
	db := controllers.GetDB()
	redisClient := controllers.GetRedisClient()
	ctx := controllers.GetContext()

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open uploaded file"})
		return
	}
	defer f.Close()

	xlFile, err := xlsx.OpenReaderAt(f, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read Excel file"})
		return
	}

	if err := validateExcelFile(xlFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sheet := xlFile.Sheets[0]
	var g errgroup.Group

	for i, row := range sheet.Rows {
		// Skip header row
		if i == 0 {
			continue
		}

		// Ensure there are enough cells in the row
		if len(row.Cells) < 10 {
			continue
		}

		// Extract data from cells
		firstName := row.Cells[0].String()
		lastName := row.Cells[1].String()
		companyName := row.Cells[2].String()
		address := row.Cells[3].String()
		city := row.Cells[4].String()
		county := row.Cells[5].String()
		postal := row.Cells[6].String()
		phone := row.Cells[7].String()
		email := row.Cells[8].String()
		web := row.Cells[9].String()

		// Use the loop index as the closure variable
		rowIndex := i

		g.Go(func() error {
			// Insert data into the database
			_, err := db.Exec(
				"INSERT INTO employees (first_name, last_name, company_name, address, city, county, postal, phone, email, web) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				firstName, lastName, companyName, address, city, county, postal, phone, email, web,
			)
			if err != nil {
				log.Printf("Error inserting data into database at row %d: %v", rowIndex, err)
				return err
			}

			// Store data in Redis Set (Allowing duplicate email IDs)
			employeeID := fmt.Sprintf("employee:%d", rowIndex) // Use a unique identifier for each employee
			redisKey := fmt.Sprintf("employees:%s", employeeID)
			if err := redisClient.SAdd(ctx, redisKey, email).Err(); err != nil {
				log.Printf("Error storing data in Redis for row %d: %v", rowIndex, err)
				// Handle error as needed
			}
			if err := redisClient.Expire(ctx, redisKey, 5*time.Minute).Err(); err != nil {
				log.Printf("Error setting expiration for Redis key for row %d: %v", rowIndex, err)
				// Handle error as needed
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "File uploaded successfully"})
}
