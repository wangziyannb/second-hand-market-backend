package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"SecondHandMarketBackend/model"
	"SecondHandMarketBackend/service"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

/**
 * @description: Place a new order.
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */
func orderPlaceHandler(w http.ResponseWriter, r *http.Request) {
	//require: order info, bundled as a json data
	//how-to: check user profile and product, then add this order to db

	//check if this user is trying to buy one of products he sells
	decoder := json.NewDecoder(r.Body)
	var order model.Order
	if err := decoder.Decode(&order); err != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}
	if order.BuyerId == order.SellerId {
		http.Error(w, "Recursive purchase", http.StatusBadRequest)
		return
	}
	//no need to check check if login Buyer exists due to foreign key constraint

	//no need to check product due to foreign key constraint
	//but we need to check seller though we have foreign key constraint
	//might be another one? Or this product might be unavailable to purchase.
	var product model.Product
	product.ID = order.ProductId
	product, err := service.SearchProductByID(&product)
	if err != nil {
		http.Error(w, "No such product", http.StatusBadRequest)
		return
	}
	if product.UserId != order.SellerId {
		http.Error(w, "Seller mismatch", http.StatusBadRequest)
		return
	}
	if product.State != "for sale" {
		http.Error(w, "unavaliable to purchase", http.StatusBadRequest)
		return
	}
	//quickly use token to get buyer info
	Buyer := service.GetUserByToken(r.Context().Value("user").(*jwt.Token).Claims)
	//or we can double check the buyer, depends on security level
	//Buyer, err:=service.CheckUserByToken(r.Context().Value("user").(*jwt.Token).Claims)
	// if err != nil {
	// 	http.Error(w, "Unknown user or user profile is broken", http.StatusUnauthorized)
	// 	return
	// }
	if order.BuyerId != Buyer.ID {
		http.Error(w, "Login user is not the buyer", http.StatusUnauthorized)
		return
	}
	//save order
	order.State = "pending"
	err = service.CreateOrder(&order)
	if err != nil {
		fmt.Fprint(w, "Failed to establish the order. Check the validation of ids for this order")
		return
	}
	fmt.Fprint(w, "New order established")
	//product states change
	service.ChangeProductState(order.ProductId, "pending")
}

/**
 * @description: order state change
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */

func orderStateChangeHandler(w http.ResponseWriter, r *http.Request) {
	//parse in orderid and state
	//four states: pending, shipping, completed, canceled

	fmt.Println("Received an order state change request")

	var order model.Order
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 0, 64)
	if err != nil {
		http.Error(w, "Failed to parse order id to uint", http.StatusInternalServerError)
		return
	}
	order.ID = uint(id)

	//check if this order exists
	if _, err = service.CheckOrderByID(&order); !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Fail to find order", http.StatusInternalServerError)
		fmt.Printf("Fail to find order %v.\n", err)
		return
	}

	//get orderId and state from request
	state := r.URL.Query().Get("state")

	//change order state to Cancelled
	if err := service.ChangeOrderState(&order, state); err != nil {
		http.Error(w, "Fail to change order state", http.StatusInternalServerError)
		fmt.Printf("Fail to change order state %v.\n", err)
		return
	}

	fmt.Fprint(w, "Order state has changed")

}

/**
 * @description: cancel order
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */
func orderCancelHandler(w http.ResponseWriter, r *http.Request) {
	//parse in token and orderId
	//extrac user email from token
	//verify user is the buyer
	//cancel order:change order state to "Cancelled"

	fmt.Println("Received an order cancel request")

	user := r.Context().Value("user") //extract user token
	claims := user.(*jwt.Token).Claims
	userEmail := claims.(jwt.MapClaims)["Email"]

	//get orderId from request
	var order model.Order
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 0, 64)
	if err != nil {
		http.Error(w, "Failed to parse order id to uint", http.StatusInternalServerError)
		return
	}
	order.ID = uint(id)

	//check if this order exists
	//result is a slice of order,i.e., there may be more than one order (or none).
	if _, err = service.CheckOrderByID(&order); !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Fail to find order", http.StatusInternalServerError)
		fmt.Printf("Fail to find order %v.\n", err)
		return
	}

	if order.Buyer.Email == userEmail {
		//change order state to cancelled
		if err := service.ChangeOrderState(&order, "Cancelled"); err != nil {
			http.Error(w, "Fail to cancel order", http.StatusInternalServerError)
			fmt.Printf("Fail to cancel order %v.\n", err)
			return
		}
	} else {
		http.Error(w, "Unauthorized to cancel order", http.StatusInternalServerError)
		fmt.Println("Unauthorized to cancel order")
	}

}
