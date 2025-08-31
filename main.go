package main

import (
	"fmt"
	"log"
	"restapi/config"
	"restapi/repository"
	"restapi/routes"
	"restapi/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	conn, err := config.Clickhouse()
	if err != nil {
		log.Fatal("❌ Error conectando a ClickHouse:", err)
	}
	log.Println("✅ Conexión establecida con ClickHouse")

	bitacoraRepo := repository.NewBitacoraRepository(conn)
	bitacoraService := services.NewBitacoraService(bitacoraRepo)

	api := app.Group("/api")
	routes.Routes(api, bitacoraService)

	host := config.Config("HOST")
	port := config.Config("PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("🚀 Aplicación escuchando en http://%s\n", addr)

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("❌ No se pudo levantar el servicio: %v", err)
	}
}
