package models

type userRole string

const (
	Admin userRole = "admin"
	Seller userRole = "seller"
	Customer userRole = "customer"
)

type User struct {
	UserId uint `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email string `gorm:"unique;not null"`
	Phone string `gorm:"unique"`
	Password string `gorm:"unique;not null"`
	Role userRole `gorm:"unique"`
}