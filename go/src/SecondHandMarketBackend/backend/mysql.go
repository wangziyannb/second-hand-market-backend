package backend

import (
	"SecondHandMarketBackend/constants"
	"SecondHandMarketBackend/model"

	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	/**
	 * @description: we can use MysqlBE to:
	 *          1、use MysqlBE.Db to manually build query (chain method)
	 *          2、use MysqlBE.Db.* to read&write
	 */
	MysqlBE *MysqlBackend
)

type MysqlBackend struct {
	Db *gorm.DB
}

/**
 * @description: Initialize Mysql backend
 * @return {*}
 */
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
	err = db.AutoMigrate(&model.Order{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Product{})
	if err != nil {
		log.Fatal(err)
	}
	// err = db.AutoMigrate(&model.Message{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = db.AutoMigrate(&model.Conversation{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	MysqlBE = &MysqlBackend{Db: db}
}

/**
 * @description: search if any table contains one object (model.*)
 * @param {interface{}} receiver, receive the result
 * @param {*gorm.DB} query, chain method build via MysqlBE.Db, but no finisher method yet
 * @return {*} error or nil
 */
func (backend *MysqlBackend) ReadOneFromMysql(receiver interface{}, query *gorm.DB) error {
	result := query.First(receiver)
	return result.Error
}

func (backend *MysqlBackend) ReadProductFromMysql(product *model.Product, query *gorm.DB) error {
	result := query.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Email", "Phone", "UserName", "University")
	}).First(product)
	return result.Error
}

/**
 * @description: save a object to table. The table is defined by struct in model.*
 * @param {interface{}} saved  the object we want to save
 * @return {*} error or nil
 */
func (backend *MysqlBackend) SaveToMysql(saved interface{}) error {
	result := backend.Db.Create(saved)
	return result.Error
}

func (backend *MysqlBackend) UpdateToMysql() {

}

func (backend *MysqlBackend) DeleteFromMysql() {

}
