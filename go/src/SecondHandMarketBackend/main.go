package main

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("started-service")
	backend.InitMysqlBackend()
	backend.InitGCSBackend()
	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}
