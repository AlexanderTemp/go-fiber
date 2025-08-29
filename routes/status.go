package routes

import "github.com/gofiber/fiber/v2"

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "200",
		"message": "Servicio ejecutandose con normalidad 🧏‍♂️",
	})
}
