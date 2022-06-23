package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"SecondHandMarketBackend/model"
	"SecondHandMarketBackend/service"

	jwt "github.com/form3tech-oss/jwt-go"
)

/**
 * @description: post product
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one upload request")

	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	email := claims.(jwt.MapClaims)["Email"]
	university := claims.(jwt.MapClaims)["University"]
	phone := claims.(jwt.MapClaims)["Phone"]
	username := claims.(jwt.MapClaims)["UserName"]

	p := model.Product{
		ProductName: r.FormValue("ProductName"),
		Price:       r.FormValue("Price"),
		Description: r.FormValue("Description"),
		University:  university.(string),
		State:       r.FormValue("State"),
		Condition:   r.FormValue("Condition"),
	}

	quantity, err := strconv.Atoi(r.FormValue("Qty"))
	if err != nil {
		http.Error(w, "Quantity cannot be parsed into int", http.StatusBadRequest)
		fmt.Printf("Quantity cannot be parsed into int %v\n", err)
		return
	}
	p.Qty = quantity

	u := model.User{
		Email:      email.(string),
		University: university.(string),
		UserName:   username.(string),
		Phone:      phone.(string),
	}
	result, err := service.CheckUser(&u)
	if err != nil {
		http.Error(w, "Couldn't find user", http.StatusBadRequest)
		fmt.Printf("Couldn't find user %v\n", err)
		return
	}
	p.UserId = result.ID

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Couldn't parse multipart form", http.StatusBadRequest)
		fmt.Printf("Couldn't parse multipart form %v\n", err)
		return
	}
	// 之后可以改进，使用goroutine waitgroup来同时上传多个文件
	var ph model.Photo
	for _, fh := range r.MultipartForm.File["Photo"] {
		file, err := fh.Open()
		if err != nil {
			http.Error(w, "Image file is not available", http.StatusBadRequest)
			fmt.Printf("Image file is not available %v\n", err)
			return
		}
		err = service.SaveProductToGCS(&ph, &p, file)
		if err != nil {
			http.Error(w, "Couldn't save post to GCS", http.StatusBadRequest)
			fmt.Printf("Couldn't save post to GCS %v\n", err)
			return
		}
		file.Close()
	}
	// Convert model.Photo instance to JSON data
	photoJSON, err := json.Marshal(ph)
	if err != nil {
		http.Error(w, "Failed to convert photo to JSON", http.StatusInternalServerError)
		fmt.Printf("Failed to convert photo to JSON %v\n", err)
		return
	}
	p.Photo = photoJSON

	err = service.SaveProductToMysql(&p)
	if err != nil {
		http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to backend %v\n", err)
		return
	}

	fmt.Println("Post is saved successfully.")
}
