package main

import (
	"fmt"
	"log"
	"restapi/config"
	"restapi/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	api := app.Group("/api")
	routes.Routes(api)

	err := app.Listen(fmt.Sprintf(":%s", config.Config("PORT")))

	if err != nil {
		log.Println("No se pudo levantar el servicio: ", err)
	}
}
