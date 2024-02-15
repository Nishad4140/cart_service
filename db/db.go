package db

import (
	"github.com/Nishad4140/cart_service/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(connect string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connect), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entities.Cart{})
	db.AutoMigrate(&entities.CartItems{})

	return db, nil
}
