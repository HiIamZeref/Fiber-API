package routes

import (
	"errors"

	"github.com/HiIamZeref/Fiber-API/database"
	"github.com/HiIamZeref/Fiber-API/models"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	// This is not the model from the database, see this as a serializer
	ID uint `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

func CreateResponseUser(userModel models.User) User {
	return User{
		ID: userModel.ID,
		FirstName: userModel.FirstName,
		LastName: userModel.LastName,
	}
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error",
			"error": err.Error(),
		})
	}

	database.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}

	database.Database.Db.Find(&users) // Get all users from the database and store them in the users variable

	// Create a slice of User structs to hold my formatted response
	responseUsers := []User{}

	for _, user := range users {
		responseUsers = append(responseUsers, CreateResponseUser(user))
		// Create a response user for each user in the users variable
	}

	return c.Status(200).JSON(responseUsers)
}

func findUser(id int, user *models.User) error {
	database.Database.Db.Find(&user, "id = ?", id)

	if user.ID == 0 {
		return errors.New("User not found")
	}
	return nil
}

func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "PLEASE PROVIDE A VALID ID",
			"error": err.Error(),
		})
	}
	var user models.User

	if err := findUser(id, &user); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
			"error": err.Error(),
		})
	}
	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)

	
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	// Check if the id is valid
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "PLEASE PROVIDE A VALID ID",
			"error": err.Error(),
		})
	}
	
	var user models.User

	// Check if the user exists
	if err := findUser(id, &user); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
			"error": err.Error(),
		})
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
	}

	var updateUser UpdateUser

	// Parse the request body into the updateUser struct
	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error",
			"error": err.Error(),
		})
	}

	// Update the user
	user.FirstName = updateUser.FirstName
	user.LastName = updateUser.LastName

	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "PLEASE PROVIDE A VALID ID",
			"error": err.Error(),
		})
	}

	var user models.User

	if err := findUser(id, &user); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
			"error": err.Error(),
		})
	}

	if err := database.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Error",
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "User deleted successfully",
		"status": "OK",
	})

	
}