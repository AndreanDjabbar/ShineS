package controllers

import (
	"fmt"
	"net/http"
	"shines/middlewares"
	"shines/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ViewShopHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	profile := models.Profile{}
	err := models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", GetuserId(c)).First(&profile).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/shop-information-page", c)
		return
	}

	shop := models.Shop{}
	err = models.DB.Model(&models.Shop{}).Select("*").Where("User_id = ?", GetuserId(c)).First(&shop).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/shop-information-page", c)
		return
	}
	context := gin.H{
		"title":       "Shop Information",
		"image":       profile.Image,
		"shopName":    shop.ShopName,
		"address":     shop.ShopAddress,
		"description": shop.ShopDescription,
		"shopImage":   shop.ShopImage,
		"isAdmin":     IsAdmin(c),
	}

	isSeller := IsSeller(c)
	if !isSeller {
		context["isSeller"] = false
		context["buttonCmnd"] = "Register"
	} else {
		context["isSeller"] = true
		context["buttonCmnd"] = "Update"
	}

	c.HTML(
		http.StatusOK,
		"profileShop.html",
		context,
	)
}

func ShopHandler(c *gin.Context) {
	shop := models.Shop{}
	profile := models.Profile{}
	userId := GetuserId(c)

	models.DB.Model(&models.Shop{}).Select("*").Where("User_id = ?", userId).First(&shop)
	models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", userId).First(&profile)

	var shopName, description, address string
	var shopNameErr, descriptionErr, addressErr, fileErr string

	shopName = c.PostForm("shopName")
	description = c.PostForm("description")
	address = c.PostForm("address")

	if len(shopName) < 5 {
		shopNameErr = "Minimum Shop Name is 5 Characters!"
	}

	if len(address) < 5 {
		addressErr = "Minimum Shop Address is 5 Characters!"
	}

	file, err := c.FormFile("photo")
	if file == nil {
		if shopNameErr == "" && descriptionErr == "" && addressErr == "" {
			shop := models.Shop{
				UserID:          uint(userId),
				ShopName:        shopName,
				ShopDescription: description,
				ShopAddress:     address,
			}
			err := models.DB.Model(&models.Shop{}).Where("user_id = ?", userId).Updates(&shop).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/shop-information-page", c)
				return
			}
			SetRole(c)
			c.Redirect(
				http.StatusFound,
				"/shines/main/shop-information-page",
			)
			return
		}
		context := gin.H{
			"title":          "Shop Information",
			"image":          profile.Image,
			"shopName":       shop.ShopName,
			"address":        shop.ShopAddress,
			"description":    shop.ShopDescription,
			"shopImage":      shop.ShopImage,
			"addressErr":     addressErr,
			"shopNameErr":    shopNameErr,
			"descriptionErr": descriptionErr,
			"isAdmin":        IsAdmin(c),
		}

		isSeller := IsSeller(c)
		if !isSeller {
			context["isSeller"] = false
			context["buttonCmnd"] = "Register"
		} else {
			context["isSeller"] = true
			context["buttonCmnd"] = "Update"
		}

		c.HTML(
			http.StatusOK,
			"profileShop.html",
			context,
		)
	} else {
		if err != nil {
			fileErr = "Failed Upload Picture"
		}
		err = c.SaveUploadedFile(file, "views/images/"+file.Filename)
		if err != nil {
			fileErr = "Failed Upload Picture"
		}
		if shopNameErr == "" && descriptionErr == "" && addressErr == "" && fileErr == "" {
			shop = models.Shop{
				UserID:          uint(userId),
				ShopName:        shopName,
				ShopDescription: description,
				ShopAddress:     address,
				ShopImage:       file.Filename,
			}
			err = models.DB.Model(&models.Shop{}).Where("user_id = ?", userId).Updates(&shop).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/shop-information-page", c)
				return
			}
			SetRole(c)
			c.Redirect(
				http.StatusFound,
				"/shines/main/shop-information-page",
			)
			return
		}
		context := gin.H{
			"title":          "Shop Information",
			"shopName":       shop.ShopName,
			"address":        shop.ShopAddress,
			"shopImage":      shop.ShopImage,
			"image":          profile.Image,
			"shopNameErr":    shopNameErr,
			"descriptionErr": descriptionErr,
			"addressErr":     addressErr,
			"fileErr":        fileErr,
			"isAdmin":        IsAdmin(c),
		}

		isSeller := IsSeller(c)
		if !isSeller {
			context["isSeller"] = false
			context["buttonCmnd"] = "Register"
		} else {
			context["isSeller"] = true
			context["buttonCmnd"] = "Update"
		}

		c.HTML(
			http.StatusOK,
			"profileShop.html",
			context,
		)
	}
}

func ViewDetailShopHandler(c *gin.Context) {
	role := GetRole(c)
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	if role != "Admin" {
		c.Redirect(
			http.StatusFound,
			"shines/main/home-page",
		)
		return
	}
	strUserId := c.Param("userId")
	userId, _ := strconv.Atoi(strUserId)
	profile := models.Profile{}
	err := models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", userId).First(&profile).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/administrator-page", c)
		return
	}

	shop := models.Shop{}
	err = models.DB.Model(&models.Shop{}).Select("*").Where("User_id = ?", userId).First(&shop).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/administrator-page", c)
		return
	}
	context := gin.H{
		"title":         "Detail Shop Information",
		"image":         profile.Image,
		"shopName":      shop.ShopName,
		"address":       shop.ShopAddress,
		"description":   shop.ShopDescription,
		"shopImage":     shop.ShopImage,
		"isSeller":      IsSeller(c),
		"isAdmin":       IsAdmin(c),
		"isAdminTarget": IsAdminTarget(c, userId),
		"userID":        userId,
	}
	c.HTML(
		http.StatusOK,
		"detailShop.html",
		context,
	)
}

func DetailShopHandler(c *gin.Context) {
	role := GetRole(c)
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	if role != "Admin" {
		c.Redirect(
			http.StatusFound,
			"shines/main/home-page",
		)
		return
	}
	strUserId := c.Param("userId")
	userId, _ := strconv.Atoi(strUserId)
	profile := models.Profile{}
	err := models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", userId).First(&profile).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/administrator-page", c)
		return
	}

	shop := models.Shop{}
	err = models.DB.Model(&models.Shop{}).Select("*").Where("User_id = ?", userId).First(&shop).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/administrator-page", c)
		return
	}

	var shopName, description, address string
	var shopNameErr, addressErr, fileErr string

	shopName = c.PostForm("shopName")
	description = c.PostForm("description")
	address = c.PostForm("address")

	if len(shopName) < 5 && len(shopName) != 0 {
		shopNameErr = "Minimum Shop Name is 5 Characters!"
	}

	if len(address) < 5 && len(address) != 0 {
		addressErr = "Minimum Shop Address is 5 Characters!"
	}

	file, err := c.FormFile("photo")
	if file == nil {
		if shopNameErr == "" && addressErr == "" {
			shop := models.Shop{
				UserID:          uint(userId),
				ShopName:        shopName,
				ShopDescription: description,
				ShopAddress:     address,
			}
			err := models.DB.Model(&models.Shop{}).Where("user_id = ?", userId).Updates(&shop).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/shop-information-page", c)
				return
			}
			SetRoleTarget(c, userId)
			targetUrl := fmt.Sprintf("/shines/main/detail-shop-page/%d", userId)
			c.Redirect(
				http.StatusFound,
				targetUrl,
			)
			return
		}
		context := gin.H{
			"title":         "Shop Information",
			"image":         profile.Image,
			"shopName":      shop.ShopName,
			"address":       shop.ShopAddress,
			"description":   shop.ShopDescription,
			"shopImage":     shop.ShopImage,
			"addressErr":    addressErr,
			"shopNameErr":   shopNameErr,
			"isAdmin":       IsAdmin(c),
			"userID":        userId,
			"isAdminTarget": IsAdminTarget(c, userId),
		}
		c.HTML(
			http.StatusOK,
			"detailShop.html",
			context,
		)
	} else {
		if err != nil {
			fileErr = "Failed Upload Picture"
		}
		err = c.SaveUploadedFile(file, "views/images/"+file.Filename)
		if err != nil {
			fileErr = "Failed Upload Picture"
		}
		if shopNameErr == "" && addressErr == "" && fileErr == "" {
			shop = models.Shop{
				UserID:          uint(userId),
				ShopName:        shopName,
				ShopDescription: description,
				ShopAddress:     address,
				ShopImage:       file.Filename,
			}
			err = models.DB.Model(&models.Shop{}).Where("user_id = ?", userId).Updates(&shop).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/shop-information-page", c)
				return
			}
			SetRoleTarget(c, userId)
			targetUrl := fmt.Sprintf("/shines/main/detail-shop-page/%d", userId)
			c.Redirect(
				http.StatusFound,
				targetUrl,
			)
			return
		}
		context := gin.H{
			"title":         "Shop Information",
			"shopName":      shop.ShopName,
			"address":       shop.ShopAddress,
			"shopImage":     shop.ShopImage,
			"image":         profile.Image,
			"shopNameErr":   shopNameErr,
			"addressErr":    addressErr,
			"fileErr":       fileErr,
			"isAdmin":       IsAdmin(c),
			"userID":        userId,
			"isAdminTarget": IsAdminTarget(c, userId),
		}
		c.HTML(
			http.StatusOK,
			"detailShop.html",
			context,
		)
	}
}