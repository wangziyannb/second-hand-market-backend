package handler

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var mySigningKey = []byte("secret")

func InitRouter() http.Handler {
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
	router.Handle("/product/{id}", jwtmiddleware.Handler(http.HandlerFunc(productHandler))).Methods("GET")
	router.Handle("/product-state-change/{id}",
		jwtmiddleware.Handler(http.HandlerFunc(productStateChangeHandler))).Methods("POST")

	// router.Handle("/search/{name}", jwtmiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("POST")
	router.Handle("/order-place", jwtmiddleware.Handler(http.HandlerFunc(orderPlaceHandler))).Methods("POST")
	// router.Handle("/order-history", jwtmiddleware.Handler(http.HandlerFunc(orderHistoryHandler))).Methods("POST")
	router.Handle("/order-detail/{id}", jwtmiddleware.Handler(http.HandlerFunc(orderDetailHandler))).Methods("POST")
	// router.Handle("/order-cancel/{id}", jwtmiddleware.Handler(http.HandlerFunc(orderCancelHandler))).Methods("POST")
	// router.Handle("/order-state-change/{id}",
	// 	jwtmiddleware.Handler(http.HandlerFunc(orderStateChangeHandler))).Methods("POST")

	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")
	router.Handle("/user-check/{id}", http.HandlerFunc(checkUserHandler)).Methods("GET")

	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

	return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}
