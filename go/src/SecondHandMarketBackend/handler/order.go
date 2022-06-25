package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"SecondHandMarketBackend/model"
	"SecondHandMarketBackend/service"

	jwt "github.com/form3tech-oss/jwt-go"
	"gorm.io/gorm"
	"github.com/gorilla/mux"
)

/**
 * @description: order state change
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */

func orderStateChangeHandler(w http.ResponseWriter, r *http.Request) {
	//parse in orderid and state
	//three states: Pending, Shipping, Completed

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

	if order.Buyer.Email == userEmail{
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
