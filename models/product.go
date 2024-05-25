package models

type Product struct {
	ProductId uint   `gorm:"primaryKey"`
	ProductName string `gorm:"size:100"`
	ProductDescription string `gorm:"size:200"`
	ProductPrice float64 `gorm:"not null"`
	ProductStock uint `gorm:"not null"`
	ProductImage string `gorm:"size:200"`
	ProductCategory string `gorm:"type:ENUM('Bunga', 'Peralatan', 'Bibit', 'Pupuk', 'Aksesoris', 'Others');not null"`
	ShopId uint
	Shop Shop `gorm:"foreignKey:ShopId"`
}