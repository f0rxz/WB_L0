package main

import (
	"context"
	"log"
	"orderservice/internal/infrastructure/connectors"

	"github.com/gofiber/fiber/v2"
)

func main() {
	ctx := context.Background()

	db, err := connectors.ConnectPostgres(ctx)
	defer db.Close()

	if err != nil {
		panic(err)
	}
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Here will be order service"))
	})
	log.Fatal(app.Listen(":3000"))
}
