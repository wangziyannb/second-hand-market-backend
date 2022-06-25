package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"
)

/**
 * @description: use a order to search if this order exists
 * @param {*model.Order} order
 * @return {*}
 */

func CheckOrderByID(order *model.Order) (model.Order, error) {
	var result model.Order
	//build query via chain method
	query := backend.MysqlBE.Db.Where(&order)
	err := backend.MysqlBE.ReadOneFromMysql(&result, query)
	return result, err
}

func ChangeOrderState(order *model.Order, state string) error {
	query := backend.MysqlBE.Db.Where(&order)
	return backend.MysqlBE.UpdateToMysql(query, state)
}

