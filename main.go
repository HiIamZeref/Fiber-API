package main

import (
	"log"

	"github.com/HiIamZeref/Fiber-API/database"
	"github.com/HiIamZeref/Fiber-API/routes"
	"github.com/gofiber/fiber/v2"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to the API!")
}

func setupRoutes(app *fiber.App){
	// Welcome endpoint
	app.Get("/api", welcome)

	// User endpoints
	app.Get("/api/users", routes.GetUsers)
	app.Post("/api/users", routes.CreateUser)
	app.Get("/api/users/:id", routes.GetUser)
	app.Put("/api/users/:id", routes.UpdateUser)
	app.Delete("/api/users/:id", routes.DeleteUser)

	// Product endpoints
	app.Get("/api/products", routes.GetProducts)
	app.Post("/api/products", routes.CreateProduct)
	app.Get(("/api/products/:id"), routes.GetProduct)
	app.Put("/api/products/:id", routes.UpdateProduct)
	app.Delete("/api/products/:id", routes.DeleteProduct)
}

func main(){
	database.ConnectDb()
	app := fiber.New()

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}