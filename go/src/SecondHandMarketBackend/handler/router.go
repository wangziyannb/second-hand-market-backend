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
	router.Handle("/message/new", jwtmiddleware.Handler(http.HandlerFunc(messageNewUploadHandler))).Methods("POST")
	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")
	
	return router
}
