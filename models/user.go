package models

type userRole string

const (
	Admin userRole = "admin"
	Seller userRole = "seller"
	Customer userRole = "customer"
)

type User struct {
    UserId   uint     `gorm:"primaryKey"`
    Username string   `gorm:"unique;size:100;not null"`
    Email    string   `gorm:"not null;size:100"`
    Phone    string   `gorm:"unique"`
    Password string   `gorm:"unique;not null"`
    Role     string `gorm:"type:ENUM('Admin', 'Seller', 'Customer');not null"`
}