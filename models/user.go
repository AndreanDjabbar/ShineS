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
    Role     userRole `gorm:"type:ENUM('Admin', 'Seller', 'Customer');not null"`
}

type Profile struct {
    ProfileId uint   `gorm:"primaryKey"`
    FirstName string `gorm:"size:50;"`
    LastName  string `gorm:"size:50;"`
    Address   string `gorm:"size:200;"`
    Image     string `gorm:"size:200"`
    UserID    uint
    User      User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type Shop struct {
    SellerId        uint   `gorm:"primaryKey"`
    ShopName        string `gorm:"size:100"`
    ShopDescription string `gorm:"size:200"`
    ShopAddress     string `gorm:"size:200"`
    ShopImage       string `gorm:"size:200"`
    UserID          uint
    User            User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}