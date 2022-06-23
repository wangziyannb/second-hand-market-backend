package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"
)

/**
 * @description: create a new user in users table
 * @param {*model.User} user
 * @return {*}
 */
func CreateUser(user *model.User) error {
	err := backend.MysqlBE.SaveToMysql(&user)
	return err
}

/**
 * @description: use a user to search if this user exists
 * @param {*model.User} user
 * @return {*}
 */
func CheckUser(user *model.User) (model.User, error) {
	var result model.User
	//build query via chain method
	query := backend.MysqlBE.Db.Where(&user)
	err := backend.MysqlBE.ReadOneFromMysql(&result, query)
	return result, err
}
