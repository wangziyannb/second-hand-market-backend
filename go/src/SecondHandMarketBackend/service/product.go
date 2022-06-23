package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"
	"mime/multipart"
	"strconv"
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

func SearchProductByID(product *model.Product) (model.Product, error) {
	var result model.Product
	//build query via chain method
	query := backend.MysqlBE.Db.Where(&product)
	err := backend.MysqlBE.ReadProductFromMysql(&result, query)
	return result, err
}
