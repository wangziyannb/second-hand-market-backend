package backend

import (
	"SecondHandMarketBackend/constants"
	"SecondHandMarketBackend/model"
	"fmt"

	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlBE *MysqlBackend
)

type MysqlBackend struct {
	Db *gorm.DB
}

func InitMysqlBackend() {
	url := constants.DB_USER + ":" + constants.DB_PWD + "@tcp(" + constants.DB_URL + ")/" + constants.DB_NAME + "?parseTime=true&loc=Local"
	//"laioffer_test:123456@tcp(212.64.40.29:3306)/laioffer_test"
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxOpenConns(5)
	sqlDb.SetMaxIdleConns(2)
	sqlDb.SetConnMaxIdleTime(time.Minute)

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Message{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Conversation{})
	if err != nil {
		log.Fatal(err)
	}
	MysqlBE = &MysqlBackend{Db: db}
}

func (backend *MysqlBackend) ReadFromMysql() {

}

func (backend *MysqlBackend) ReadOneFromMysql(receiver interface{}, query *gorm.DB) error{
	result := query.First(receiver)
	return result.Error
}

func (backend *MysqlBackend) SaveToMysql(i interface{}) error {
	fmt.Printf("%T", i)
	result := backend.Db.Create(i)
	return result.Error
}

func (backend *MysqlBackend) UpdateToMysql() {

}

func (backend *MysqlBackend) DeleteFromMysql() {

}
