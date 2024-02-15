package initializer

import (
	"github.com/Nishad4140/cart_service/adapter"
	"github.com/Nishad4140/cart_service/service"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB) *service.CartService {
	adapter := adapter.NewCartAdapter(db)
	service := service.NewCartService(adapter)

	return service
}
