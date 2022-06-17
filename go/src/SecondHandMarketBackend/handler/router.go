package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/message/new", http.HandlerFunc(messageNewUploadHandler)).Methods("POST")
	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin",http.HandlerFunc(signinHandler)).Methods("POST")
	return router
}
