package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"
	"fmt"
	"mime/multipart"
	"strconv"
)

func SaveProductToGCS(photo *model.Photo, product *model.Product, file multipart.File) error {

	// Generate unique name for each photo
	uniqueName := strconv.FormatUint(uint64(product.ID), 10) + product.ProductName + strconv.Itoa(len(photo.Photos))
	fmt.Print(uniqueName)

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
