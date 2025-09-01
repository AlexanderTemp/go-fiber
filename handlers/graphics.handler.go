package handlers

import (
	"log"
	"restapi/dto"
	"restapi/services"

	"github.com/gofiber/fiber/v2"
)

func TendenciaConsumo(c *fiber.Ctx, service *services.BitacoraService) error {
	var filtros dto.FiltroTendenciaConsumoDto
	if err := c.QueryParser(&filtros); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"mensaje": "Parámetros inválidos",
		})
	}

	res, err := service.ObtenerTendenciaConsumo(c.Context(), filtros)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"mensaje": "Error al obtener el registro de tendencias de consumo",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"mensaje": "Registros de tendencia en consumo obtenidos",
		"data":    res,
	})
}

func ObtenerDatosExtra(c *fiber.Ctx, service *services.BitacoraService) error {
	var filtros dto.FiltroExtraDto
	if err := c.QueryParser(&filtros); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"mensaje": "Parámetros inválidos",
		})
	}

	res, err := service.ObtenerDatosExtra(c.Context(), filtros)
	if err != nil {
		log.Fatal(err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"mensaje": "Error al obtener los datos extra",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"mensaje": "Datos extra obtenidos",
		"data":    res,
	})
}
