package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger(c *fiber.Ctx) error {
	start := time.Now()

	// Continue stack
	err := c.Next()

	// Get the response status and body
	status := c.Response().StatusCode()

	duration := time.Since(start)

	// Log request and response details
	log.Printf("[%s] %s %s %s %d %s", start.Format(time.RFC3339), c.Method(), c.Path(), c.IP(), status, duration)

	return err
}
