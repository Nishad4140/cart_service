package adapter

import (
	"fmt"

	"github.com/Nishad4140/cart_service/entities"
	"gorm.io/gorm"
)

type CartAdapter struct {
	DB *gorm.DB
}

func NewCartAdapter(db *gorm.DB) *CartAdapter {
	return &CartAdapter{
		DB: db,
	}
}

func (cart *CartAdapter) CreateCart(userId uint) error {
	fmt.Println("user id", userId)
	query := "INSERT INTO carts (user_id) VALUES ($1)"

	if err := cart.DB.Exec(query, userId).Error; err != nil {
		return err
	}
	return nil
}

func (cart *CartAdapter) AddToCart(req entities.CartItems, userId uint) error {
	tx := cart.DB.Begin()

	var cartId int
	var current entities.CartItems

	queryId := "SELECT id FROM carts WHERE user_id = ?"
	if err := tx.Raw(queryId, userId).Scan(&cartId).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("cart not found")
	}
	queryCurrent := "SELECT * FROM cart_items WHERE cart_id = $1 AND product_id = $2"
	if err := tx.Raw(queryCurrent, cartId, req.ProductId).Scan(&current).Error; err != nil {
		tx.Rollback()
		return err
	}
	var res entities.CartItems
	if current.ProductId == 0 {
		insertQuery := "INSERT INTO cart_items (cart_id, product_id, quantity, total) VALUES ($1, $2, $3, 0) RETURNING id, product_id, cart_id"
		if err := tx.Raw(insertQuery, cartId, req.ProductId, req.Quantity).Scan(&res).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		updateQuery := "UPDATE cart_items SET quantity = quantity + $1 WHERE cart_id = $2 AND product_id = $3"
		if err := tx.Exec(updateQuery, req.Quantity, cartId, req.ProductId).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	updateTotal := `UPDATE cart_items SET total = total + $1 WHERE cart_id = $2 AND product_id = $3`
	if err := tx.Exec(updateTotal, (req.Total * float64(req.Quantity)), cartId, req.ProductId).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
