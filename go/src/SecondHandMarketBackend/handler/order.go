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
 * @description: change order state 
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */

func orderStateChangeHandler(w http.ResponseWriter, r *http.Request) {
	//use url to get ordre Id
	//decode json data to get state
	//four states: pending, shipping, completed, canceled
	//orderStateChangeHandler can change order state to pending, shipped and completed
	//order canel is controlled by orderCancelHandler

	//Note: Seller can modify order state from pending to shipped or completed
	//Maybe can use token to extract user info and verify user is seller


	fmt.Println("Received an order state change request")

	//get order from url
	var order model.Order
	orderId, err := strconv.ParseUint(mux.Vars(r)["id"], 0, 64)
	if err != nil {
		http.Error(w, "Failed to parse order id to uint", http.StatusInternalServerError)
		return
	}
	order.ID = uint(int(orderId))

	//check if this order exists
	order, err = service.CheckOrder(&order)
	if err != nil {
		http.Error(w, "Fail to find order", http.StatusInternalServerError)
		fmt.Printf("Fail to find order %v.\n", err)
		return
	}

	//get state decode json
	//state is extracted by decoding json data
	//Note: state is passed in as Order struct with only State filled
	decoder := json.NewDecoder(r.Body)
	var neworder model.Order
	if err := decoder.Decode(&neworder); err != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}
	//check if state is valid
	if neworder.State == "pending" || neworder.State == "shipped" || neworder.State == "completed" {
		//change order state
		if err := service.ChangeOrderState(&order, neworder.State); err != nil {
			http.Error(w, "Fail to change order state", http.StatusInternalServerError)
			fmt.Printf("Fail to change order state %v.\n", err)
			fmt.Println("order state", order.State)
			return
		}
	}else{
		http.Error(w, "Invalid order state", http.StatusInternalServerError)
		fmt.Printf("Invalid order state %v.\n", err)
		fmt.Println("neworder state", neworder.State)
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
	//get user info from token
	//get orderId from url
	//verify user is the buyer
	//order can be canceled only if order is pending
	//cancel order:change order state to "canceled"
	
	//Note: we are not eagerly deleting the order 
	//User can retrieve their canceled order later
	//maybe need another function to delete canceled order if canceled order is 5-year old?

	fmt.Println("Received an order cancel request")

	//get user info from token
	var user model.User
	token := r.Context().Value("user") //extract user token
	claims := token.(*jwt.Token).Claims
	userId := claims.(jwt.MapClaims)["ID"].(float64)
	userEmail := claims.(jwt.MapClaims)["Email"].(string)
	user.ID = uint(int(userId))
	user.Email = userEmail

	//check if user exits
	user, err := service.CheckUser(&user)
	if err != nil {
		http.Error(w, "Fail to find user ", http.StatusBadRequest)
		return
	}

	//get orderId from url
	var order model.Order
	orderId, err := strconv.ParseUint(mux.Vars(r)["id"], 0, 64)
	if err != nil {
		http.Error(w, "Failed to parse order id to uint", http.StatusInternalServerError)
		return
	}
	order.ID = uint(int(orderId))

	//check if this order exists
	order, err = service.CheckOrder(&order)
	if err != nil {
		http.Error(w, "Fail to find order", http.StatusInternalServerError)
		fmt.Printf("Fail to find order %v.\n", err)
		return
	}

	//check if user is buyer
	if order.BuyerId != user.ID {
		http.Error(w, "Unauthorized to cancel order", http.StatusInternalServerError)
		fmt.Println("Unauthorized to cancel order")
		return
	}

	//check if order is pending
	if order.State != "pending" {
		http.Error(w, "Order has shipped and fail to cancel", http.StatusInternalServerError)
		fmt.Println("Order has shipped and fail to cancel")
		return
	}

	//cancel order -> change order state to canceled
	if err := service.ChangeOrderState(&order, "canceled"); err != nil {
		http.Error(w, "Fail to change order state", http.StatusInternalServerError)
		fmt.Printf("Fail to change order state %v.\n", err)
		return
	}

	fmt.Fprint(w, "Order has cancelled")

}
