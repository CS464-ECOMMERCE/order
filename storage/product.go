package storage

import (
	"order/models"

	"gorm.io/gorm"
)

type ProductInterface interface {
	RevertProductQuantity(id, quantity uint64, tx *gorm.DB) error
}

type ProductDB struct {
	write *gorm.DB
}

func NewProductTable(write *gorm.DB) ProductInterface {
	StorageInstance.AutoMigrate(&models.Product{})
	return &ProductDB{
		write: write,
	}
}

func (p *ProductDB) RevertProductQuantity(id, quantity uint64, tx *gorm.DB) error {
	Product := &models.Product{}
	db := tx
	if db == nil {
		db = p.write
	}

	ret := db.Where("id = ?", id).First(&Product)
	if ret.Error != nil {
		return ret.Error
	}

	newQuantity := Product.Inventory + quantity
	ret = db.Model(&models.Product{}).Where("id = ?", id).Update("inventory", newQuantity)
	if ret.Error != nil {
		return ret.Error
	}

	return nil
}
