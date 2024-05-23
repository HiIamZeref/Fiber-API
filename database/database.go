package database

import (
	"log"
	"os"

	"github.com/HiIamZeref/Fiber-API/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database")
		os.Exit(2)
	}

	log.Println("Connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running migrations")
	// TODO: Add migrations
	db.AutoMigrate(&models.Order{}, &models.Product{}, &models.Order{})

	Database = DbInstance{Db: db}
}