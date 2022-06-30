package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"

)

/**
 * @description: check if this order exists
 * @param {*model.Order} order
 * @return {model.order, *}
 */
func CheckOrder(order *model.Order) (model.Order, error) {
	var result model.Order
	//build query via chain method
	query := backend.MysqlBE.Db.Where(&order)
	err := backend.MysqlBE.ReadOneFromMysql(&result, query)
	return result, err
}

/**
 * @description: change order state
 * @param {*model.Order} order, state string
 * @return {*}
 */
func ChangeOrderState(order *model.Order, state string) error {
	query := backend.MysqlBE.Db.Model(&order)
	return backend.MysqlBE.UpdateOneToMysql(query, "State", state)
}

/**
 * @description: create a new order
 * @param {*model.Order} order
 * @return {*}
 */
func CreateOrder(order *model.Order) error {
	err := backend.MysqlBE.SaveToMysql(&order)
	return err
}

