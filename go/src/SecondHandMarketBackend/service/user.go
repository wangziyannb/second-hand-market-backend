package service

import (
	"SecondHandMarketBackend/backend"
	"SecondHandMarketBackend/model"

	"github.com/form3tech-oss/jwt-go"
	"gorm.io/gorm"
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

/**
 * @description: use token to quickly gain the user info via token without checking the validtaion
 * @param {jwt.Claims} claims
 * @return {model.User}
 *
 * @e.g. : user := service.GetUserByToken(r.Context().Value("user").(*jwt.Token).Claims)
 */
func GetUserByToken(claims jwt.Claims) model.User {
	u := model.User{
		Email:      claims.(jwt.MapClaims)["Email"].(string),
		University: claims.(jwt.MapClaims)["University"].(string),
		UserName:   claims.(jwt.MapClaims)["UserName"].(string),
		Phone:      claims.(jwt.MapClaims)["Phone"].(string),
	}
	id := claims.(jwt.MapClaims)["ID"].(float64)
	u.ID = uint(int(id))
	return u
}

/**
 * @description: use token to quickly check the validation of user. Normally we do not have to do so.
 * @param {jwt.Claims} claims
 * @return {model.User, error} the model.User is a user built by token
 *
 * @e.g. : user, err := service.CheckUserByToken(r.Context().Value("user").(*jwt.Token).Claims)
 */
func CheckUserByToken(claims jwt.Claims) (model.User, error) {
	u := model.User{
		Email:      claims.(jwt.MapClaims)["Email"].(string),
		University: claims.(jwt.MapClaims)["University"].(string),
		UserName:   claims.(jwt.MapClaims)["UserName"].(string),
		Phone:      claims.(jwt.MapClaims)["Phone"].(string),
	}
	v, err := CheckUser(&u)
	if err != nil {
		return u, err
	}
	id := claims.(jwt.MapClaims)["ID"].(float64)
	if uint64(int(id)) == uint64(v.ID) {
		return v, nil
	} else {
		return u, gorm.ErrRecordNotFound
	}
}

/**
 * @description: use ID to quickly check the validation of user
 * @param {uint} ID
 * @return {model.User, error} the model.User is a user built by token
 */
func CheckUserByID(ID uint) (model.User, error) {
	var u model.User
	u.ID = ID
	v, err := CheckUser(&u)
	if err != nil {
		return u, err
	} else {
		return v, nil
	}
}
