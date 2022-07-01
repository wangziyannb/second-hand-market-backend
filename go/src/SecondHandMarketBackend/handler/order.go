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
	//no need to check if login Buyer exists due to foreign key constraint

	//no need to check product due to foreign key constraint
	//but we need to check seller though we have foreign key constraint
	//might be another one? Or this product might be unavailable to purchase.
	product, err := service.SearchProductByID(order.ProductId)
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
	//use url to get ordre Id->check order using ID, and get old order state
	//decode json data to get new order state
	//use token to get user, and verify if user is associated with order
	//change order state: ->seller can change order state from pending to shipped
	//               ->buyer can change state from shipped(via mail)/pending (in-person) to completed   
	//if order is completed, change product state to "hidden"

	fmt.Println("Received an order state change request")

	//get orderId
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 0, 64)
	if err != nil {
		http.Error(w, "Failed to parse order id to uint", http.StatusInternalServerError)
		return
	}

	//check if this order exists using orderID
	order, err := service.CheckOrderByID(uint(id))
	if err != nil {
		http.Error(w, "Fail to find order", http.StatusBadRequest)
		fmt.Printf("Fail to find order %v.\n", err)
		return
	}

	oldstate := order.State

	//check user via token
	token := r.Context().Value("user") //extract user token
	claims := token.(*jwt.Token).Claims
	user,err := service.CheckUserByToken(claims) 

	if err!=nil {
		http.Error(w, "User not exists", http.StatusUnauthorized)
		fmt.Printf("User not exists %v.\n", err)
		return
	}

	//check if user is associated with the order
	if order.BuyerId != user.ID && order.SellerId != user.ID {
		http.Error(w, "Unauthorised user", http.StatusUnauthorized)
		fmt.Printf("Unauthorised user %v.\n", err)
		return
	}

	//get state decode json
	//state is extracted by decoding json data
	//Note: state is passed in as Order struct with only State filled
	decoder := json.NewDecoder(r.Body)
	var neworder model.Order
	if err := decoder.Decode(&neworder); err != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		fmt.Printf("Bad json %v.\n", err)
		return
	}

	//check if user, old state and newstate is valid
	//if seller, pending->shipped
	//if buyer, shipped->completed
	if oldstate == "pending" && order.SellerId == user.ID && neworder.State == "shipped" || oldstate == "shipped" && order.BuyerId == user.ID && neworder.State == "completed" {
		//change order state
		if err := service.ChangeOrderState(order.ID, neworder.State); err != nil {
			http.Error(w, "Fail to change order state", http.StatusInternalServerError)
			fmt.Printf("Fail to change order state %v.\n", err)
			return
		}
		//change product state
		if neworder.State == "completed" {
			service.ChangeProductState(order.ProductId, "hidden")
		}
	}else{
		http.Error(w, "Fail to change order state", http.StatusUnauthorized)
		fmt.Printf("Fail to change order state")
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
	//verify if user is the buyer
	//order can be canceled only if order state is pending
	//cancel order:change order state to "canceled"
	//change product state to "for sale"
	
	//Note: we are not eagerly deleting the order 
	//User can retrieve their canceled order later
	//maybe need another function to delete canceled order if canceled order is 5-year old?

	fmt.Println("Received an order cancel request")

	//get user info from token
	token := r.Context().Value("user") //extract user token
	claims := token.(*jwt.Token).Claims
	user, err := service.CheckUserByToken(claims)

	if err != nil {
		http.Error(w, "Fail to find user ", http.StatusBadRequest)
		return
	}

	//get orderId
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 0, 64)
	if err != nil {
		http.Error(w, "Failed to parse order id to uint", http.StatusInternalServerError)
		return
	}

	//check if this order exists using orderID
	order, err := service.CheckOrderByID(uint(id))
	if err != nil {
		http.Error(w, "Fail to find order", http.StatusBadRequest)
		fmt.Printf("Fail to find order %v.\n", err)
		return
	}

	oldstate := order.State

	//check if user is buyer
	if order.BuyerId != user.ID {
		http.Error(w, "Unauthorized to cancel order", http.StatusUnauthorized)
		fmt.Println("Unauthorized to cancel order")
		return
	}

	//check if order is pending
	if oldstate != "pending"{
		fmt.Println("order can't be canceled")
		return
	}
	
	//cancel order -> change order state to canceled
	if err := service.ChangeOrderState(order.ID, "canceled"); err != nil {
		http.Error(w, "Fail to change order state", http.StatusInternalServerError)
		fmt.Printf("Fail to change order state %v.\n", err)
		return
	}

	fmt.Fprint(w, "Order has cancelled")

	//change product state to "for sale"
	service.ChangeProductState(order.ProductId, "for sale")

}

/**
 * @description: order detail
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */
func orderDetailHandler(w http.ResponseWriter, r *http.Request) {
	//check order

	//Note: service.checkOrderById has changed. Related code below has been commented out/modified accordingly.

	fmt.Println("Recevied one item order detail request")
	//var order model.Order
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 0, 64)
	if err != nil {
		http.Error(w, "Failed to parse order id to unit", http.StatusInternalServerError)
		return
	}
	//order.ID = uint(id)

	order, err := service.CheckOrderByID(uint(id))

	if err != nil {
		http.Error(w, "No such order", http.StatusBadRequest)
		return
	}
	//check user validation
	//token -> id -> SellerId or BuyerId
	user := service.GetUserByToken(r.Context().Value("user").(*jwt.Token).Claims)
	if user.ID == order.BuyerId || user.ID == order.SellerId {
		//jsonify
		js, err := json.Marshal(order)
		if err != nil {
			http.Error(w, "Failed to get json data from search result", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		http.Error(w, "Failed to get order detail request", http.StatusBadRequest)
		return
	}
}

func orderHistoryHandler(w http.ResponseWriter, r *http.Request) {
	//use token to get user
	//return all orders of the user, regardless if user is a buyer and user

	//get user from token
	token := r.Context().Value("user") //extract user token
	claims := token.(*jwt.Token).Claims
	user, err := service.CheckUserByToken(claims) 

	if err!=nil {
		http.Error(w, "User not exists", http.StatusUnauthorized)
		fmt.Printf("User not exists %v.\n", err)
		return
	}

	//search order by user
	
	orders, err := service.SearchOrderByUser(user.ID)
	if err!=nil{
		http.Error(w, "Fail to find order", http.StatusInternalServerError)
		fmt.Printf("Fail to find orders %v.\n", err)
		return
	}

	js, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "Failed to get json data from search result", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func unit(i int) {
	panic("unimplemented")
}