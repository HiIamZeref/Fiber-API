package routes

import (
	"time"

	"github.com/HiIamZeref/Fiber-API/database"
	"github.com/HiIamZeref/Fiber-API/models"
	"github.com/gofiber/fiber/v2"
)

type Order struct{
	ID 		uint `json:"id" gorm:"primaryKey"`
	User 	User `json:"user"`
	Product Product `json:"product"`
	CreatedAt time.Time `json:"order_date"`
}

// CreateResponseOrder creates a response order
func CreateResponseOrder(order models.Order, user User, product Product) Order {
	return Order{
		ID: order.ID,
		User: user,
		Product: product,
		CreatedAt: order.CreatedAt,
	}
}

// CreateOrder creates an order in the database
func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	// Parsing the request body and storing it in the order variable
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error",
			"error": err.Error(),
		})
	}

	// Checking if the user exists
	var user models.User
	if err := findUser(order.UserRefer, &user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error",
			"error": err.Error(),
		})
	}

	// Checking if product exits
	var product models.Product
	if err := findProduct(order.ProductRefer, &product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error",
			"error": err.Error(),
		})
	}

	// Create the order in the database
	database.Database.Db.Create(&order)

	// Format the response
	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	// Return the created order
	return c.Status(200).JSON(responseOrder)
																

}


//GetOrders returns all orders in the database
func GetOrders(c *fiber.Ctx) error {
	// Create a slice of orders
	orders := []models.Order{}

	// Get all orders from the database
	database.Database.Db.Find(&orders)

	// Create a slice of Order structs to hold the formatted response
	responseOrders := []Order{}

	// Iterate over the orders and create a response order for each
	for _, order := range orders {
		// Instantiate the user and product
		var user models.User
		var product models.Product

		// Find the user and product for the order
		database.Database.Db.Find(&user, "id = ?", order.UserRefer)
		database.Database.Db.Find(&product, "id = ?", order.ProductRefer)

		// Create the response user and product first for the order
		responseUser := CreateResponseUser(user)
		responseProduct := CreateResponseProduct(product)

		// Create the response order
		responseOrders = append(responseOrders, CreateResponseOrder(order, responseUser, responseProduct))

	}

	// Return the response orders
	return c.Status(200).JSON(responseOrders)
}

// Helper function to an order by ID
func FindOrder(id int, order *models.Order) error {
	// Find the order in the database
	database.Database.Db.Find(&order, "id = ?", id)

	// Check if the order exists
	if order.ID == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Order not found")
	}

	return nil
}

// GetOrder returns a single order by ID
func GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var order models.Order

	// Check if ID is valid
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "PLEASE PROVIDE A VALID ID",
			"error":   err.Error(),
		})
	}

	// Check if the order exists
	if err := FindOrder(id, &order); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error",
			"error": err.Error(),
		})
	}

	// Instantiate the user and product
	var user models.User
	var product models.Product

	// Get user and product from the database
	database.Database.Db.First(&user, "id = ?", order.UserRefer)
	database.Database.Db.First(&product, "id = ?", order.ProductRefer)

	// Create the response user and product
	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)

	// Create the response order
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	// Return the response order
	return c.Status(200).JSON(responseOrder)
}