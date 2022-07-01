package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"
	"errors"
)

func CheckOrderByID(ID uint) (model.Order, error) {
	var order model.Order
	order.ID = ID
	result, err := CheckOrder(&order)
	if err != nil {
		return order, err
	} else {
		return result, nil
	}
}

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
func ChangeOrderState(ID uint, newState string) error {
	var order model.Order
	order.ID = ID
	query := backend.MysqlBE.Db.Model(&order)
	switch newState {
	case "pending",
		"shipped",
		"completed",
		"canceled":
		return backend.MysqlBE.UpdateOneToMysql(query, "state", newState)
	}
	return errors.New("not a valid state")
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

func SearchOrderByUser(ID uint) ([]model.Order, error) {
	var orders []model.Order

	//build query via chain method
	query := backend.MysqlBE.Db.Where(backend.MysqlBE.Db.Table("Order").Where("buyer_id", ID)).Or(backend.MysqlBE.Db.Table("Order").Where("seller_id", ID))
	err := backend.MysqlBE.ReadAllFromMysql(&orders, query)

	return orders, err
}
