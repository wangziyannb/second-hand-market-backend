package handler

import (
	"SecondHandMarketBackend/model"
	"SecondHandMarketBackend/service"
	"encoding/json"
	"fmt"
	"net/http"
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
	err := service.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to add user to backend", http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "New user profile established")
}
