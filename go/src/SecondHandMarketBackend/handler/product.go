package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"SecondHandMarketBackend/model"
	"SecondHandMarketBackend/service"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
)

/**
 * @description: post product
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one upload request")

	user := service.GetUserByToken(r.Context().Value("user").(*jwt.Token).Claims)

	p := model.Product{
		ProductName: r.FormValue("ProductName"),
		Price:       r.FormValue("Price"),
		Description: r.FormValue("Description"),
		University:  user.University,
		State:       "for sale",
		Condition:   r.FormValue("Condition"),
		UserId:      user.ID,
	}

	quantity, err := strconv.Atoi(r.FormValue("Qty"))
	if err != nil {
		http.Error(w, "Quantity cannot be parsed into int", http.StatusBadRequest)
		return
	}
	p.Qty = quantity

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Couldn't parse multipart form", http.StatusBadRequest)
		return
	}
	// 之后可以改进，使用goroutine waitgroup来同时上传多个文件
	var ph model.Photo
	for _, fh := range r.MultipartForm.File["Photo"] {
		file, err := fh.Open()
		if err != nil {
			http.Error(w, "Image file is not available", http.StatusBadRequest)
			return
		}
		err = service.SaveProductToGCS(&ph, &p, file)
		if err != nil {
			http.Error(w, "Couldn't save post to GCS", http.StatusBadRequest)
			return
		}
		file.Close()
	}
	// Convert model.Photo instance to JSON data
	photoJSON, err := json.Marshal(ph)
	if err != nil {
		http.Error(w, "Failed to convert photo to JSON", http.StatusInternalServerError)
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

/**
 * @description: return the detail information of the product based on its ID.
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */
func productHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one item detail request")

	id, err := strconv.ParseUint(mux.Vars(r)["id"], 0, 64)
	if err != nil {
		http.Error(w, "Failed to parse product id to uint64", http.StatusInternalServerError)
		return
	}

	product, err := service.SearchProductByID(uint(id))
	//bugfix by Ziyan Wang: unhandled error
	if err != nil {
		http.Error(w, "No such product", http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(product)
	if err != nil {
		http.Error(w, "Failed to get json data from search result", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/**
 * @description: change the state of product. Only seller can do this operation.
 * see also: service.ChangeProductState
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */
func productStateChangeHandler(w http.ResponseWriter, r *http.Request) {
	user := service.GetUserByToken(r.Context().Value("user").(*jwt.Token).Claims)
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 0, 64)
	if err != nil {
		http.Error(w, "Failed to parse product id to uint64", http.StatusInternalServerError)
		return
	}
	//check permission
	product, err := service.SearchProductByID(uint(id))
	if err != nil {
		http.Error(w, "No such product", http.StatusBadRequest)
		return
	}
	if product.UserId != user.ID {
		http.Error(w, "No permission to do that", http.StatusBadRequest)
		return
	}
	//get new state
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		fmt.Print(err)
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}
	//must not be pending
	switch product.State {
	case "hidden",
		"for sale":
		err = service.ChangeProductState(uint(id), product.State)
		if err != nil {
			http.Error(w, "Failed to change state of product", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, "Successfully changed the state of product")
		return
	}
	http.Error(w, "Not a valid state", http.StatusBadRequest)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one search request")
	user := service.GetUserByToken(r.Context().Value("user").(*jwt.Token).Claims)

	decoder := json.NewDecoder(r.Body)
	var product model.Product
	if err := decoder.Decode(&product); err != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}

	products, err := service.SearchProductByName(product.ProductName, user.University)

	if err != nil {
        http.Error(w, "Failed to get products from backend", http.StatusInternalServerError) //StatusInternalServerError = 500
        fmt.Printf("Failed to get products from backend %v.\n", err)
        return
    }
	js, err := json.Marshal(products)
    //handle parse过程中出现的error
    if err != nil {
        http.Error(w, "Failed to parse products into JSON format", http.StatusInternalServerError)
        fmt.Printf("Failed to parse products into JSON format %v.\n", err)
        return
    }
	w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}