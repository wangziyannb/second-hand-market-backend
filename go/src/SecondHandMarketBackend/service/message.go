package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"
)

func CreateConversation(conversation *model.Conversation) error {
	err := backend.MysqlBE.SaveToMysql(conversation)
	return err
}
