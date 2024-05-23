package routes

import (
	"github.com/HiIamZeref/Fiber-API/database"
	"github.com/HiIamZeref/Fiber-API/models"
	"github.com/gofiber/fiber/v2"
)

type Product struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func CreateResponseProduct(productModel models.Product) Product {
	return Product{
		ID:           productModel.ID,
		Name:         productModel.Name,
		SerialNumber: productModel.SerialNumber,
	}
}

// CreateProduct creates a product in the database
func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	// Parse the body of the request and store it in the product variable
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error",
			"error":   err.Error(),
		})
	}

	// Create the product in the database
	database.Database.Db.Create(&product)

	// Return the created product
	return c.Status(200).JSON(CreateResponseProduct(product))

	
}

// GetProducts returns all products in the database
func GetProducts(c *fiber.Ctx) error {
	// Initialize a slice of Product structs to hold the formatted response
	products := []models.Product{}

	// Get all products from the database
	database.Database.Db.Find(&products)

	// Create a slice of Product structs to hold the formatted response
	responseProducts := []Product{}

	// Iterate over the products and create a response product for each
	for _, product := range products {
		responseProducts = append(responseProducts, CreateResponseProduct(product))
	}

	// Return the response products
	return c.Status(200).JSON(responseProducts)
}

// Helper function to find a product by ID
func findProduct(id int, product *models.Product) error {
	database.Database.Db.Find(&product, "id = ?", id)

	if product.ID == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	return nil
}

// GetProduct returns a single product by ID
func GetProduct(c *fiber.Ctx) error {
	// Get the ID from the URL
	id, err := c.ParamsInt("id")

	// Check if the ID is valid
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "PLEASE PROVIDE A VALID ID",
			"error":   err.Error(),
		})
	}

	// Initialize a product variable
	var product models.Product

	// Find the product by ID
	if err := findProduct(id, &product); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Product not found",
			"error":   err.Error(),
		})
	}

	// Return the product
	return c.Status(200).JSON(CreateResponseProduct(product))
}

func UpdateProduct(c *fiber.Ctx) error {
	// Get the ID from the URL
	id, err := c.ParamsInt("id")

	// Check if the ID is valid
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "PLEASE PROVIDE A VALID ID",
			"error":   err.Error(),
		})
	}
	var product models.Product

	// Find the product by ID
	if err := findProduct(id, &product); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Product not found",
			"error":   err.Error(),
		})
	}
	type UpdateProduct struct {
		Name 	   string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}

	// Initialize a product variable
	var updateProduct UpdateProduct
	
	// Parse the body of the request and store it in the updateProduct variable
	if err := c.BodyParser(&updateProduct); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error",
			"error":   err.Error(),
		})
	}

	// Update the product with the new values
	product.Name = updateProduct.Name
	product.SerialNumber = updateProduct.SerialNumber

	// Save the updated product to the database
	database.Database.Db.Save(&product)

	return c.Status(200).JSON(CreateResponseProduct(product))
}

func DeleteProduct(c *fiber.Ctx) error {
	// Get the ID from the URL
	id, err := c.ParamsInt("id")

	// Check if the ID is valid
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "PLEASE PROVIDE A VALID ID",
			"error":   err.Error(),
		})
	}

	// Initialize a product variable
	var product models.Product

	// Find the product by ID
	if err := findProduct(id, &product); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Product not found",
			"error":   err.Error(),
		})
	}

	if err := database.Database.Db.Delete(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Product deleted",
		"status": "OK",
	})
}