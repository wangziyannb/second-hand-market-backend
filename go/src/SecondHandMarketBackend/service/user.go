package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"
	"fmt"
)

func CreateUser(user *model.User) error {
	err := backend.MysqlBE.SaveToMysql(&user)
	return err
}

func CheckUser(user *model.User) bool {
	var a model.User
	fmt.Print("var:", a)
	err := backend.MysqlBE.ReadFromMysql(&user, &a)
	if err != nil {
		return false
	}
	fmt.Print(a)
	return true
}
