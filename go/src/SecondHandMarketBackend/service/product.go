package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"
	"errors"
	"mime/multipart"
	"strconv"

	"gorm.io/gorm"
)

func SaveProductToGCS(photo *model.Photo, product *model.Product, file multipart.File) error {
	// Generate unique name for each photo
	uniqueName := strconv.FormatUint(uint64(product.ID), 10) + product.ProductName + strconv.Itoa(len(photo.Photos))

	imagelink, err := backend.GCSBackend.SaveToGCS(file, uniqueName)
	if err != nil {
		return err
	}
	photo.Photos = append(photo.Photos, imagelink)
	return nil
}

func SaveProductToMysql(product *model.Product) error {
	return backend.MysqlBE.SaveToMysql(product)
}

func SearchProductByID(ID uint) (model.Product, error) {
	var product model.Product
	product.ID = ID
	//build query via chain method
	query := backend.MysqlBE.Db.Where(&product).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Email", "Phone", "UserName", "University")
	})
	err := backend.MysqlBE.ReadOneFromMysql(&product, query)
	return product, err
}

func ChangeProductState(ID uint, newState string) error {
	var product model.Product
	product.ID = ID
	query := backend.MysqlBE.Db.Model(&product)
	switch newState {
	case "hidden",
		"pending",
		"for sale":
		return backend.MysqlBE.UpdateOneToMysql(query, "state", newState)
	}
	return errors.New("not a valid state")
}

func SearchProductByName(productName string, university string) ([]model.Product, error) {

	var products []model.Product
	query := backend.MysqlBE.Db.Where("product_name LIKE ? AND university = ? AND state = ?", "%" + productName + "%", university, "for sale");

	err := backend.MysqlBE.ReadAllFromMysql(&products, query)
	if err != nil {
            return nil, err
    }

    return products, err
}