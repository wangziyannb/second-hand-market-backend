/*
 * @Author: xyzhao009 79874305+xyzhao009@users.noreply.github.com
 * @Date: 2022-06-30 13:16:21
 * @LastEditors: xyzhao009 79874305+xyzhao009@users.noreply.github.com
 * @LastEditTime: 2022-06-30 18:15:09
 * @FilePath: /second-hand-market-backend-3/go/src/SecondHandMarketBackend/handler/router.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
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
	router.Handle("/order-history", jwtmiddleware.Handler(http.HandlerFunc(orderHistoryHandler))).Methods("GET")

	// router.Handle("/order-detail/{id}", jwtmiddleware.Handler(http.HandlerFunc(orderDetailHandler))).Methods("POST")
	router.Handle("/order-cancel/{id}", jwtmiddleware.Handler(http.HandlerFunc(orderCancelHandler))).Methods("POST")
	router.Handle("/order-state-change/{id}",
		jwtmiddleware.Handler(http.HandlerFunc(orderStateChangeHandler))).Methods("POST")

	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")
	router.Handle("/user-check/{id}", http.HandlerFunc(checkUserHandler)).Methods("GET")

	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

	return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}
