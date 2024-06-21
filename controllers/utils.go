package controllers

import (
	"fmt"
	"net/http"
	"shines/middlewares"
	"shines/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func IsSeller(c *gin.Context) bool {
	role := GetRole(c)
	if role == "Seller" {
		return true
	}
	return false
}

func Add1(x int) int {
	return x + 1
}

func IsAdmin(c *gin.Context) bool {
	role := GetRole(c)
	if role == "Admin" {
		return true
	}
	return false
}

func GetuserId(c *gin.Context) int {
	user := middlewares.GetSession(c)
	var userId int
	models.DB.Model(&models.User{}).Select("UserId").Where("username = ?", user).First(&userId)
	return userId
}

func GetSellerId(c *gin.Context) int {
	userId := GetuserId(c)
	var sellerId int
	models.DB.Model(&models.Shop{}).Select("SellerId").Where("user_id = ?", userId).First(&sellerId)
	return sellerId
}

func GetEmailUser(c *gin.Context) string {
	user := middlewares.GetSession(c)
	var email string
	models.DB.Model(&models.User{}).Select("Email").Where("username = ?", user).First(&email)
	return email
}

func GetPasswordUser(c *gin.Context) string {
	user := middlewares.GetSession(c)
	var password string
	models.DB.Model(&models.User{}).Select("Password").Where("username = ?", user).First(&password)
	return password
}

func CreateProfile(c *gin.Context) {
	userId := GetuserId(c)

	var count int64
	models.DB.Model(&models.Profile{}).Where("user_id = ?", userId).Count(&count)

	if count > 0 {
		return
	} else {
		profile := models.Profile{
			UserID: uint(userId),
			Image:  "default.png",
		}
		err := models.DB.Create(&profile).Error
		if err != nil {

			ErrorHandler1("Failed to Create Data", "/shines/main/personal-information-page", c)
			return
		}
	}
}

func CreateShop(c *gin.Context) {
	userId := GetuserId(c)
	var count int64
	models.DB.Model(&models.Shop{}).Where("user_id = ?", userId).Count(&count)

	if count > 0 {
		return
	} else {
		shop := models.Shop{
			UserID:    uint(userId),
			ShopImage: "store.png",
		}
		err := models.DB.Create(&shop).Error
		if err != nil {

			ErrorHandler1("Failed to Create Data", "/shines/main/personal-information-page", c)
			return
		}
	}
}

func AddToCart(c *gin.Context, productID int, quantity int, stock int) {
	userId := GetuserId(c)
	urlSource := fmt.Sprintf("/shines/main/detail-product-page/%d", productID)
	cart := models.Cart{}
	err := models.DB.Model(&models.Cart{}).Where("user_id = ? AND product_id = ?", userId, productID).First(&cart).Error
	if err != nil {
		cart.UserID = uint(userId)
		cart.ProductID = uint(productID)
		cart.Quantity = uint(quantity)
		err = models.DB.Create(&cart).Error
		if err != nil {

			ErrorHandler1("Failed to Create Data", urlSource, c)
			return
		}
	} else {
		newQuantity := cart.Quantity + uint(quantity)
		if newQuantity >= uint(stock) {
			cart.Quantity = uint(stock)
		}
		err = models.DB.Save(&cart).Error
		if err != nil {

			ErrorHandler1("Failed to Update Data", urlSource, c)
			return
		}
	}
}

func UpdateCart(c *gin.Context, cartID, quantity, stock int) {
	userId := GetuserId(c)
	urlSource := fmt.Sprintf("/shines/main/update-cart-page/%d", cartID)
	cart := models.Cart{}
	err := models.DB.Model(&models.Cart{}).Where("user_id = ? AND cart_id = ?", userId, cartID).First(&cart).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", urlSource, c)
		return
	}
	newQuantity := uint(quantity)
	fmt.Println(newQuantity)
	if newQuantity >= uint(stock) {
		cart.Quantity = uint(stock)
	} else {
		cart.Quantity = newQuantity
	}
	err = models.DB.Save(&cart).Error
	if err != nil {

		ErrorHandler1("Failed to Update Data", urlSource, c)
		return
	}
}

func DeleteCart(c *gin.Context, cartID int) {
	userId := GetuserId(c)
	urlSource := fmt.Sprintf("/shines/main/cart-page")
	cart := models.Cart{}
	err := models.DB.Model(&models.Cart{}).Where("user_id = ? AND cart_id = ?", userId, cartID).First(&cart).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", urlSource, c)
		return
	}
	err = models.DB.Delete(&cart).Error
	if err != nil {

		ErrorHandler1("Failed to Delete Data", urlSource, c)
		return
	}
}

func GetNameProduct(c *gin.Context, productId int) string {
	product := models.Product{}
	models.DB.Model(&models.Product{}).Select("product_name").Where("product_id = ?", productId).First(&product)
	return product.ProductName
}

func AddToTransaction(c *gin.Context, price float64, productID int, quantityOrder int) {
	userID := GetuserId(c)
	transaction := models.Transactions{}
	transaction.UserID = uint(userID)
	transaction.ProductPrice = price
	transaction.Quantity = uint(quantityOrder)
	transaction.ProductName = GetNameProduct(c, productID)
	now := time.Now()
	transaction.TransactionDate = now.Format("2006-01-02")
	transaction.ProductID = uint(productID)
	err := models.DB.Create(&transaction).Error
	if err != nil {

		ErrorHandler1("Failed to Create Data", "/shines/main/cart-page", c)
		return
	}
}

func UpdateStockProduct(c *gin.Context, productId int, quantityOrder int) {
	product := models.Product{}
	models.DB.Model(&models.Product{}).Where("product_id = ?", productId).First(&product)
	product.ProductStock = product.ProductStock - uint(quantityOrder)
	err := models.DB.Save(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Update Data", "/shines/main/cart-page", c)
		return
	}
}

func ClearCart(c *gin.Context) {
	userId := GetuserId(c)
	cart := models.Cart{}
	err := models.DB.Model(&models.Cart{}).Where("user_id = ?", userId).Delete(&cart).Error
	if err != nil {

		ErrorHandler1("Failed to Delete Data", "/shines/main/cart-page", c)
		return
	}
}

func GetRoleTarget(c *gin.Context, userId int) string {
	user := models.User{}
	models.DB.Model(&models.User{}).Select("*").Where("User_id = ?", userId).First(&user)
	return string(user.Role)
}

func GetPriceProduct(c *gin.Context, productId int) float64 {
	product := models.Product{}
	models.DB.Model(&models.Product{}).Select("product_price").Where("product_id = ?", productId).First(&product)
	return product.ProductPrice
}

func IsAdminTarget(c *gin.Context, userId int) bool {
	role := GetRoleTarget(c, userId)
	if role == "Admin" {
		return true
	}
	return false
}

func GetRole(c *gin.Context) string {
	userId := GetuserId(c)
	user := models.User{}

	models.DB.Model(&models.User{}).Select("*").Where("User_id = ?", userId).First(&user)
	return string(user.Role)
}

func SetRole(c *gin.Context) {
	userId := GetuserId(c)
	user := models.User{}
	models.DB.Model(&models.User{}).Select("*").Where("User_id = ?", userId).First(&user)
	currentRole := GetRole(c)
	if currentRole == "Customer" {
		user.Role = "Seller"
		err := models.DB.Model(&models.User{}).Where("user_id = ?", userId).Updates(&user).Error
		if err != nil {

			ErrorHandler1("Failed to Update Data", "/shines/main/shop-information-page", c)
			return
		}
		return
	} else {
		return
	}
}

func SetRoleTarget(c *gin.Context, userId int) {
	user := models.User{}
	models.DB.Model(&models.User{}).Select("*").Where("User_id = ?", userId).First(&user)
	currentRole := GetRole(c)
	if currentRole == "Customer" {
		user.Role = "Seller"
		err := models.DB.Model(&models.User{}).Where("user_id = ?", userId).Updates(&user).Error
		if err != nil {

			ErrorHandler1("Failed to Update Data", "/shines/main/shop-information-page", c)
			return
		}
		return
	} else {
		return
	}
}

func GetShopId(c *gin.Context) int {
	userId := GetuserId(c)
	var shopId int
	err := models.DB.Model(&models.Shop{}).Select("seller_id").Where("user_id = ?", userId).First(&shopId).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/shop-information-page", c)
		return 0
	}
	return shopId
}

func DeleteProduct(c *gin.Context, productId string) {
	var product models.Product
	models.DB.Where("product_id = ?", productId).First(&product)
	err := models.DB.Delete(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Delete Data", "/shines/main/seller-catalog-page", c)
		return
	}
	c.Redirect(http.StatusFound, "/shines/main/seller-catalog-page")
}

func isNumber(strings string) bool {
	for a := 0; a < len(strings); a++ {
		_, err := strconv.Atoi(string(strings[a]))
		if err != nil {
			return false
		}
	}
	return true
}