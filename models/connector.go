package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	database, err := gorm.Open(
		mysql.Open(
			"Kelompok3:Kelompok3Shines@tcp(localhost:3306)/shines",
		),
	)
	if err != nil {
		panic(err)
	}
	database.AutoMigrate(&User{}, &Profile{}, &Shop{}, &Product{}, &Cart{}, &Transactions{})
	DB = database
}
