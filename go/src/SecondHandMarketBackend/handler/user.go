package handler

import (
	"SecondHandMarketBackend/model"
	"SecondHandMarketBackend/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"gorm.io/gorm"
)

/**
 * @description: generate a new user
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */
func signupHandler(w http.ResponseWriter, r *http.Request) {
	//known attrs: Email, UserName, UserPwd, Phone, University
	/*need to do:
	1:make a new user and orm automatically update the table `users`
	2:return success to front end
	*/
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if user.Email == "" || user.UserPwd == "" || user.University == "" || user.Phone == "" || user.UserName == "" {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	//check if this email already in the db
	if err := service.CheckUser(&model.User{Email: user.Email}); !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Failed to add user to backend", http.StatusBadRequest)
		return
	}
	if err := service.CreateUser(&user); err != nil {
		http.Error(w, "Backend DB reports error", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "New user profile established")
}

/**
 * @description: sign in
 * @param {http.ResponseWriter} w
 * @param {*http.Request} r
 * @return {*}
 */
func signinHandler(w http.ResponseWriter, r *http.Request) {
	//check if we can log in
	//known attrs: Email, UserPwd
	/*need to do:
	1:check "users" table to find out that if this user exists via gorm
	2:return token to the front end
	*/
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if err := service.CheckUser(&user); err != nil {
		http.Error(w, "user not exists, or password is wrong", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email":      user.Email,
		"University": user.University,
		"Phone":      user.Phone,
		"UserName":   user.UserName,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	//return token to the front end
	w.Write([]byte(ss))
}
