package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db}
}

func (p *ProductRepository) AddProduct(product model.Product) error {
	// return nil // TODO: replace this
	return p.db.Create(&product).Error
}

func (p *ProductRepository) ReadProducts() ([]model.Product, error) {
	// return []model.Product{}, nil // TODO: replace this
	products := []model.Product{}
	err := p.db.Find(&products).Error
	if err != nil {
		return []model.Product{}, err
	}
	return products, nil
}

func (p *ProductRepository) DeleteProduct(id uint) error {
	// return nil // TODO: replace this
	return p.db.Where("id = ?", id).Delete(&model.Product{}).Error
}

func (p *ProductRepository) UpdateProduct(id uint, product model.Product) error {
	// return nil // TODO: replace this
	return p.db.Model(&model.Product{}).Where("id = ?", id).Updates(product).Error
}
