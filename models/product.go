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

type Cart struct {
	CartId uint `gorm:"primaryKey"`
    UserID    uint `gorm:"not null"`
    User      User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
    ProductID    uint `gorm:"not null"`
    Product      Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Quantity uint `gorm:"not null"`
}

type TransactionDetail struct {
	CartID      uint
	UserID      uint
	Username    string
	Email       string
	ProductID   uint
	ProductName string
	Price       float64
	Quantity    int
}