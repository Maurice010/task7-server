package database

import (
	"log"
	"server/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Nie można połączyć z bazą danych:", err)
	}

	err = db.AutoMigrate(
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
	)
	if err != nil {
		log.Fatal("Błąd podczas migracji:", err)
	}

	DB = db

	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count == 0 {
		db.Create(&models.Product{Name: "Laptop", Price: 3000})
		db.Create(&models.Product{Name: "Telefon", Price: 1500})
	}
}
