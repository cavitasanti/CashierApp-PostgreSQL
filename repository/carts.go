package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return CartRepository{db}
}

func (c *CartRepository) ReadCart() ([]model.JoinCart, error) {
	// return []model.JoinCart{}, nil // TODO: replace this
	joinCarts := []model.JoinCart{}
	a := c.db.Table("carts").Select("carts.id, products.id as product_id, products.name, carts.quantity, carts.total_price").
		Joins("join products on carts.product_id = products.id").Scan(&joinCarts)
	if a.Error != nil {
		return []model.JoinCart{}, a.Error
	}

	return joinCarts, nil
}

func (c *CartRepository) AddCart(product model.Product) error {
	// return nil // TODO: replace this
	if c.db.Where("product_id = ?", product.ID).First(&model.Cart{}).RowsAffected == 1 {
		// if product.id == model.Cart.product_id {
		c.db.Transaction(func(tx *gorm.DB) error {
			total := product.Price - (product.Price * (product.Discount / 100))
			// update field quantity and total_price
			tx.Model(&model.Cart{}).Where("product_id = ?", product.ID).Update("quantity", gorm.Expr("quantity + ?", 1))
			tx.Model(&model.Cart{}).Where("product_id = ?", product.ID).Update("total_price", gorm.Expr("total_price + ?", total))
			// s(model.Cart{TotalPrice: gorm.Expr("total_price + ?", total)})
			// TotalPrice: gorm.Expr("total_price + ?", total)})
			// (model.Cart{Quantity: +1, TotalPrice: total})
			// (model.Cart{Quantity: +1, TotalPrice: model.Cart.Quantity * (product.Price * product.Discount)})

			// update field stock
			tx.Transaction(func(tx *gorm.DB) error {
				tx.Model(&model.Product{}).Where("id = ?", product.ID).Update("stock", gorm.Expr("stock - ?", 1))
				// tx.Model(&model.Product{}).Where("id = ?", product.id).Updates(model.Product{Stock: product.Stock - model.Cart.Quantity})
				return nil
			})
			return nil
		})
	} else {
		c.db.Transaction(func(tx *gorm.DB) error {
			// insert new data

			tx.Create(&model.Cart{ProductID: product.ID, Quantity: 1, TotalPrice: product.Price - (product.Price * (product.Discount / 100))})

			// update field stock
			tx.Transaction(func(tx *gorm.DB) error {
				tx.Model(&model.Product{}).Where("id = ?", product.ID).Update("stock", gorm.Expr("stock - ?", 1))
				// tx.Model(&model.Product{}).Where("id = ?", product.id).Updates(model.Product{Stock: product.Stock - 1})
				return nil
			})
			return nil
		})
	}
	return nil
}

func (c *CartRepository) DeleteCart(id uint, productID uint) error {
	// return nil // TODO: replace this
	c.db.Transaction(func(tx *gorm.DB) error {
		// c.db.Raw("SELECT id, name, age FROM users WHERE name = ?", "Dito").Scan(&result)

		// delete cart
		errr := tx.Where("id = ?", id).Delete(&model.Cart{}).Error
		if errr != nil {
			return errr
		}

		// update product stock
		err := tx.Model(&model.Product{}).Where("id = ?", productID).Update("stock", gorm.Expr("stock + ?", 1)).Error
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (c *CartRepository) UpdateCart(id uint, cart model.Cart) error {
	// return nil // TODO: replace this
	err := c.db.Model(&model.Cart{}).Where("id = ?", id).Updates(cart).Error
	if err != nil {
		return err
	}
	return nil
}
