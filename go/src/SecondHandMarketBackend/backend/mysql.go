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
	db *gorm.DB
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
	MysqlBE = &MysqlBackend{db: db}
}

func (backend *MysqlBackend) ReadFromMysql(i, j interface{}) error {
	// print type
	result := backend.db.Where(i).First(j)
	return result.Error
}

func (backend *MysqlBackend) SaveToMysql(i interface{}) error {
	fmt.Printf("%T", i)
	result := backend.db.Create(i)
	return result.Error
}

func (backend *MysqlBackend) UpdateToMysql() {

}

func (backend *MysqlBackend) DeleteFromMysql() {

}
