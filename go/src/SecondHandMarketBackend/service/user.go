package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"
)

func CreateUser(user *model.User) error {
	err := backend.MysqlBE.SaveToMysql(&user)
	return err
}

func CheckUser(user *model.User) error{
	var result model.User
	err := backend.MysqlBE.ReadOneFromMysql(&result, backend.MysqlBE.Db.Where(&user))
	return err
}
