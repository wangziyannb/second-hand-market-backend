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
	fmt.Print(user)
	err := service.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to add user to backend", http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "New user profile established")
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signin request")
	w.Header().Set("Content-Type", "text/plain")

	if r.Method == "OPTIONS" {
		return
	}

	//  Get User information from client
	decoder := json.NewDecoder(r.Body)
	// 直接用这个json 的decoder

	// 声明一个user的对象
	var user model.User
	// 然后去decoder这个request，并与user对象对照
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client %v\n", err)
		return
	}

	// 拿到的user对象去服务器验证
	exists := service.CheckUser(&user)
	fmt.Print("exist", exists)
	// if err != nil {
	// 	http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
	// 	fmt.Printf("Failed to read user from Elasticsearch %v\n", err)
	// 	return
	// }

	// if !exists {
	// 	http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
	// 	fmt.Printf("User doesn't exists or wrong password\n")
	// 	return
	// }

	// // 这个是什么？这个就是那个登录验证，存在本地用来验证的，有效时间24h
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"username": user.Username,
	// 	"exp":      time.Now().Add(time.Hour * 24).Unix(),
	// })
	// // 不放密码的原因是避免被反解
	// // Unix：the number of seconds elapsed since January 1, 1970 UTC

	// tokenString, err := token.SignedString(mySigningKey)
	// if err != nil {
	// 	http.Error(w, "Failed to generate token", http.StatusInternalServerError)
	// 	fmt.Printf("Failed to generate token %v\n", err)
	// 	return
	// }

	// w.Write([]byte(tokenString))
}
