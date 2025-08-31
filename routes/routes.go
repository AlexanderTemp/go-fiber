package routes

import (
	"restapi/handlers"
	"restapi/services"

	"github.com/gofiber/fiber/v2"
)

func Routes(router fiber.Router, bitacoraService *services.BitacoraService) {
	router.Get("/estado", handlers.HealthCheck)

	// Gráficas
	router.Get("/tendencia", func(c *fiber.Ctx) error {
		return handlers.TendenciaConsumo(c, bitacoraService)
	})
}
