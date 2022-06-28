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
//semantics bug: CheckOrderByID doesn't use ID(uint) as input
func CheckOrderByID(order *model.Order) (model.Order, error) {
	var result model.Order
	//build query via chain method
	query := backend.MysqlBE.Db.Where(&order)
	err := backend.MysqlBE.ReadOneFromMysql(&result, query)
	return result, err
}

func ChangeOrderState(order *model.Order, state string) error {
	query := backend.MysqlBE.Db.Where(&order)
	return backend.MysqlBE.UpdateOneToMysql(query, "state", state)
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

