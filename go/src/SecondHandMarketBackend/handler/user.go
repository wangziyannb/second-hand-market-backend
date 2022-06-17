package handler

import (
	"SecondHandMarketBackend/model"
	"SecondHandMarketBackend/service"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
)

func signupHandler(w http.ResponseWriter, r *http.Request) {
	//new user
	//known attrs: username, pwd
	/*need to do:
	1:make a new user and orm automatically update the table `users`
	2:return id to the front end
	*/
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
	}
	fmt.Print(user)
	err := service.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to add user to backend", http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "New user profile established")
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	//check if we can log in
	//known attrs: username, pwd
	/*need to do:
	1:check "users" table to find out that if this user exists via gorm
	2:return token to the front end
	*/
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
	}
	fmt.Print(user)
	err := service.CheckUser(&user)
	if err != nil {
		http.Error(w, "user not exists, or password is wrong", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.UserName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(ss))
}
