package controllers

import (
	"net/http"
	"shines/middlewares"
	"shines/models"

	"github.com/gin-gonic/gin"
)

func IsSeller(c *gin.Context) bool {
	role := GetRole(c)
	if role == "Seller" {
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
		profile := models.Profile {
			UserID: uint(userId),
			Image: "default.png",
		}
		err := models.DB.Create(&profile).Error
		if err != nil {
			context := gin.H {
				"title":"Error",
				"message":"Failed to Create Data",
				"source":"/shines/main/personal-information-page",
			}
			c.HTML(
				http.StatusInternalServerError,
				"error.html",
				context,
			)
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
		shop := models.Shop {
			UserID: uint(userId),
			ShopImage: "store.png",
		}
		err := models.DB.Create(&shop).Error
		if err != nil {
			context := gin.H {
				"title":"Error",
				"message":"Failed to Create Data",
				"source":"/shines/main/personal-information-page",
			}
			c.HTML(
				http.StatusInternalServerError,
				"error.html",
				context,
			)
			return
		}	
	}
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
				context := gin.H{
					"title":   "Error",
					"message": "Failed to Update Data",
					"source":  "/shines/main/shop-information-page",
				}
				c.HTML(
					http.StatusInternalServerError,
					"error.html",
					context,
				)
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
		context := gin.H{
			"title":   "Error",
			"message": "Failed to Get Data",
			"source":  "/shines/main/shop-information-page",
		}
		c.HTML(
			http.StatusInternalServerError,
			"error.html",
			context,
		)
		return 0
	}
	return shopId
}

func DeleteProduct(c *gin.Context, productId string) {
	var product models.Product
	models.DB.Where("product_id = ?", productId).First(&product)
	err := models.DB.Delete(&product).Error
	if err != nil {
		context := gin.H{
			"title":   "Error",
			"message": "Failed to Delete Data",
			"source":  "/shines/main/seller-catalog-page",
		}
		c.HTML(
			http.StatusInternalServerError,
			"error.html",
			context,
		)
		return
	}
	c.Redirect(http.StatusFound, "/shines/main/seller-catalog-page")
}