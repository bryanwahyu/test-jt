package server

import (
	"fmt"
	"log"
	"strconv"
	"math/rand"
	"github.com/bryanwahyu/test-jt/internal/config"
	"github.com/bryanwahyu/test-jt/internal/db"
	"github.com/bryanwahyu/test-jt/internal/phone"
	"github.com/gofiber/fiber/v2"
)

// generateRandomPhoneNumber generates a random phone number for demonstration purposes
func generateRandomPhoneNumber() string {
	// You may use a more sophisticated phone number generation mechanism in a real-world scenario
	return fmt.Sprintf("123456789%d", rand.Intn(1000))
}
// StartServer starts the Fiber server
func StartServer() error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize the database
	if err := db.InitDB(cfg); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Initialize the phone service with the database
	phoneService := &phone.Service{
		DB: db.GetDB(),
	}

	// Define routes
	app.Get("/phone", func(c *fiber.Ctx) error {
		page := c.Query("page", "1")     // Page number
		pageSize := c.Query("pageSize", "10") // Number of items per page
	
		pageNum, err := strconv.Atoi(page)
		if err != nil || pageNum < 1 {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid page number")
		}
	
		pageSizeNum, err := strconv.Atoi(pageSize)
		if err != nil || pageSizeNum < 1 {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid page size")
		}

		// Retrieve paginated phone numbers
		response, err := phoneService.GetPaginatedPhones(pageNum, pageSizeNum)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve phone numbers")
		}
	
		return c.JSON(response)
	})

	app.Post("/phone", func(c *fiber.Ctx) error {
		// Parse request body for the new phone number
		newPhone := new(phone.Phone)
		if err := c.BodyParser(newPhone); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		// Add the new phone number to the database
		if err := phoneService.AddPhone(newPhone); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to add phone number")
		}

		return c.SendString("Phone number added successfully")
	})

	app.Put("/phone/:id", func(c *fiber.Ctx) error {
		// Get the phone ID from the URL parameters
		id := c.Params("id")

		// Parse request body for the updated phone number
		updatedPhone := new(phone.Phone)
		if err := c.BodyParser(updatedPhone); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		// Update the phone number in the database
		if err := phoneService.UpdatePhone(id, updatedPhone); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to update phone number")
		}

		return c.SendString("Phone number updated successfully")
	})

	app.Delete("/phone/:id", func(c *fiber.Ctx) error {
		// Get the phone ID from the URL parameters
		id := c.Params("id")

		// Delete the phone number from the database
		if err := phoneService.DeletePhone(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete phone number")
		}

		return c.SendString("Phone number deleted successfully")
	})
	app.Get("/generate", func(c *fiber.Ctx) error {
		numToGenerate := c.Query("numToGenerate")

		// Set a default value of 200 if numToGenerate is not provided or invalid
		count := 200
		if numToGenerate != "" {
			val, err := strconv.Atoi(numToGenerate)
			if err == nil && val > 0 {
				count = val
			}
		}
	
		// Create a slice to store generated phone numbers
		phones := make([]phone.Phone, count)
	
		// Use a wait group to wait for all goroutines to finish
		var wg sync.WaitGroup
		wg.Add(count)
	
		// Generate and add phone numbers concurrently
		for i := 0; i < count; i++ {
			go func(index int) {
				defer wg.Done()
	
				phones[index].Number = generateRandomPhoneNumber()
			}(i)
		}
	
		// Wait for all goroutines to finish
		wg.Wait()
	
		// Start a transaction
		tx := phoneService.DB.Begin()
	
		// Batch insert phone numbers
		if err := tx.Create(&phones).Error; err != nil {
			// Rollback the transaction on error
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to add phone numbers")
		}
	
		// Commit the transaction
		tx.Commit()
		return c.SendString(fmt.Sprintf("%d phone numbers added successfully", count))
	})
	app.Get("/",func(c *fiber.Ctx) error{
		return c.SendString("welcome")

	})
	app.Get("/migrate", func(c *fiber.Ctx) error {
		// Run database migration
		if err := db.Migrate(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to migrate database")
		}
	
		return c.SendString("Database migration successful")
	})	
	// Start the Fiber app
	port := cfg.Port
	log.Printf("Server started on :%d", port)
	return app.Listen(fmt.Sprintf(":%d", port))
}
