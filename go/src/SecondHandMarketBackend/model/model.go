package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email      string `gorm:"not null" json:"Email"`
	UserName   string `gorm:"not null" json:"UserName"`
	UserPwd    string `gorm:"not null" json:"UserPwd"`
	Phone      string `gorm:"not null" json:"Phone"`
	University string `gorm:"not null" json:"University"`
}

type Product struct {
	gorm.Model
	ProductName string         `gorm:"not null" json:"ProductName"`
	Price       string         `gorm:"not null" json:"Price"`
	Description string         `gorm:"not null" json:"Description"`
	University  string         `gorm:"not null" json:"University"`
	State       string         `gorm:"not null" json:"State"`
	Condition   string         `gorm:"not null" json:"Condition"`
	Photo       datatypes.JSON `gorm:"not null" json:"Photo"`
	Qty         int            `gorm:"not null" json:"Qty"`
	UserId      uint
	User        User
}

type Order struct {
	gorm.Model
	SellerId        int
	Seller          User
	BuyerId         int
	Buyer           User
	ProductId       int
	Product         Product
	Qty             int            `gorm:"not null" json:"Qty"`
	PlaceTime       datatypes.Date `gorm:"not null" json:"PlaceTime"`
	FinishTime      datatypes.Date `gorm:"not null" json:"FinishTime"`
	Price           string         `gorm:"not null" json:"Price"`
	State           string         `gorm:"not null" json:"State"` //3 states, pending, shipping, complete
	DeliveryAddress string         `gorm:"not null" json:"DeliveryAddress"`
	DeliveryType    string         `gorm:"not null" json:"DeliveryType"`
}

//below is p1 to p2
type Favorites struct {
	// Product
	//
}

type Message struct {
	gorm.Model
	ConversationID int
	Message        string `gorm:"not null"`
	SenderId       int    `gorm:"not null"`
	ReceiverId     int    `gorm:"not null"`
}

type Conversation struct {
	gorm.Model
	User1Id  int
	User2Id  int
	User1    User
	User2    User
	Messages []Message
}

type Photo struct {
	Photos []string `json:"photos"`
}
