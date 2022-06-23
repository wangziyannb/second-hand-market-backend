package handler

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
)

var mySigningKey = []byte("secret")

func InitRouter() *mux.Router {
	jwtmiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	router := mux.NewRouter()
	//currently unavailable
	router.Handle("/message/new", jwtmiddleware.Handler(http.HandlerFunc(messageNewUploadHandler))).Methods("POST")

	router.Handle("/upload", jwtmiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
	// router.Handle("/product/{id}", jwtmiddleware.Handler(http.HandlerFunc(productHandler))).Methods("POST")
	// router.Handle("/product-state-change/{id}",
	// 	jwtmiddleware.Handler(http.HandlerFunc(productStateChangeHandler))).Methods("POST")

	// router.Handle("/search/{name}", jwtmiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("POST")
	// router.Handle("/order-history", jwtmiddleware.Handler(http.HandlerFunc(orderHistoryHandler))).Methods("POST")
	// router.Handle("/order-detail/{id}", jwtmiddleware.Handler(http.HandlerFunc(orderDetailHandler))).Methods("POST")
	// router.Handle("/order-cancel/{id}", jwtmiddleware.Handler(http.HandlerFunc(orderCancelHandler))).Methods("POST")
	// router.Handle("/order-state-change/{id}",
	// 	jwtmiddleware.Handler(http.HandlerFunc(orderStateChangeHandler))).Methods("POST")

	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

	return router
}
