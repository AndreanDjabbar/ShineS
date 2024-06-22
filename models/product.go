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
    BuyerID    uint `gorm:"not null"`
    Buyer      User   `gorm:"foreignKey:BuyerID;constraint:OnDelete:CASCADE;"`
	SellerID	uint `gorm:"not null"`
	Seller		User `gorm:"foreignKey:SellerID;constraint:OnDelete:CASCADE;"`
    ProductID    uint `gorm:"not null"`
    Product      Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Quantity uint `gorm:"not null"`
}

type Transactions struct {
	TransactionID uint `gorm:"primaryKey"`
	BuyerID uint `gorm:"not null"`
	Buyer User `gorm:"foreignKey:BuyerID;constraint:OnDelete:CASCADE;"`
	SellerID uint `gorm:"not null"`
	Seller User `gorm:"foreignKey:SellerID;constraint:OnDelete:CASCADE;"`
	ProductPrice float64 `gorm:"not null"`
	TransactionDate string `gorm:"not null"`
	ProductID uint `gorm:"not null"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Quantity uint `gorm:"not null"`
	ProductName string `gorm:"size:100 not null"`
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