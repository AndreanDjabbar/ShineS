package controllers

import (
	"fmt"
	"math/big"
	"net/http"
	"shines/middlewares"
	"shines/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RootHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if isLogged {
		c.Redirect(
			http.StatusFound,
			"/shines/main/home-page",
		)
	} else {
		c.Redirect(
			http.StatusFound,
			"/shines/main/login-page",
		)
	}
}

func ViewHomeHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	userID := GetuserId(c)
	products := []models.Product{}
	err := models.DB.Model(&models.Product{}).
    Select("*").
    Where("Shop_id != ? AND product_stock != ?", userID, 0).
    Find(&products).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return
	}
	context := gin.H{
		"title":    "Home",
		"products": products,
		"isSeller": IsSeller(c),
		"isAdmin":  IsAdmin(c),
	}
	fmt.Println(products)
	c.HTML(
		http.StatusOK,
		"home.html",
		context,
	)
}

func ViewPersonalHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	profile := models.Profile{}
	models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", GetuserId(c)).First(&profile)

	context := gin.H{
		"title":     "Personal Information",
		"image":     profile.Image,
		"firstName": profile.FirstName,
		"lastName":  profile.LastName,
		"address":   profile.Address,
		"isSeller":  IsSeller(c),
		"isAdmin":   IsAdmin(c),
	}
	c.HTML(
		http.StatusOK,
		"profilePersonal.html",
		context,
	)
}

func PersonalHandler(c *gin.Context) {
	userId := GetuserId(c)
	profile := models.Profile{}
	models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", userId).First(&profile)
	var firstName, lastName, address string
	var firstNameErr, lastNameErr, addressErr, fileErr string

	firstName = c.PostForm("firstname")
	lastName = c.PostForm("lastname")
	address = c.PostForm("address")

	if len(firstName) < 2 {
		firstNameErr = "Minimum First Name is 2 Character!"
	}
	if len(lastName) < 3 {
		lastNameErr = "Minimum Last Name is 3 Characters!"
	}

	if len(address) < 5 {
		addressErr = "Minimum Address is 5 Characters!"
	}

	file, err := c.FormFile("picture")
	if file == nil {
		if firstNameErr == "" && lastNameErr == "" && addressErr == "" {
			profile := models.Profile{
				UserID:    uint(userId),
				FirstName: firstName,
				LastName:  lastName,
				Address:   address,
			}
			err := models.DB.Model(&models.Profile{}).Where("user_id = ?", userId).Updates(&profile).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/personal-information-page", c)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/shines/main/personal-information-page",
			)
			return
		}
		context := gin.H{
			"title":        "Personal Information",
			"firstName":    profile.FirstName,
			"lastName":     profile.LastName,
			"address":      profile.Address,
			"firstNameErr": firstNameErr,
			"lastNameErr":  lastNameErr,
			"image":        profile.Image,
			"addressErr":   addressErr,
		}
		c.HTML(
			http.StatusOK,
			"profilePersonal.html",
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
		if firstNameErr == "" && lastNameErr == "" && addressErr == "" && fileErr == "" {
			profile := models.Profile{
				UserID:    uint(userId),
				FirstName: firstName,
				LastName:  lastName,
				Address:   address,
				Image:     file.Filename,
			}
			err := models.DB.Model(&models.Profile{}).Where("user_id = ?", userId).Updates(&profile).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/personal-information-page", c)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/shines/main/personal-information-page",
			)
			return
		}
		context := gin.H{
			"title":        "Personal Information",
			"firstName":    profile.FirstName,
			"lastName":     profile.LastName,
			"address":      profile.Address,
			"image":        profile.Image,
			"firstNameErr": firstNameErr,
			"lastNameErr":  lastNameErr,
			"addressErr":   addressErr,
			"fileErr":      fileErr,
			"isSeller":     IsSeller(c),
			"isAdmin":      IsAdmin(c),
		}
		c.HTML(
			http.StatusOK,
			"profilePersonal.html",
			context,
		)
	}
}

func ViewCredentialHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	email := GetEmailUser(c)
	profile := models.Profile{}
	err := models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", GetuserId(c)).First(&profile).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/credential-information-page", c)
		return
	}
	username := middlewares.GetSession(c)
	context := gin.H{
		"title":     "Credential Information",
		"username":  username,
		"image":     profile.Image,
		"firstName": profile.FirstName,
		"lastName":  profile.LastName,
		"address":   profile.Address,
		"email":     email,
		"isSeller":  IsSeller(c),
		"isAdmin":   IsAdmin(c),
	}

	c.HTML(
		http.StatusOK,
		"profileCredential.html",
		context,
	)
}

func CredentialHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	user := models.User{}
	err := models.DB.Model(&models.User{}).Select("*").Where("User_id = ?", GetuserId(c)).First(&user).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/credential-information-page", c)
		return
	}
	profile := models.Profile{}
	err = models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", GetuserId(c)).First(&profile).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/credential-information-page", c)
		return
	}
	var username, email, password1, password2 string
	var usernameErr, emailErr, password1Err, password2Err string

	userId := GetuserId(c)

	username = c.PostForm("username")
	email = c.PostForm("email")
	password1 = c.PostForm("password")
	password2 = c.PostForm("password-confirmation")

	if len(username) < 5 {
		usernameErr = "Minimum Username is 5 Characters!"
	}
	if len(email) < 5 {
		emailErr = "Minimum Email is 5 Characters!"
	}
	if password1 != "" && len(password1) < 5 {
		password1Err = "Minimum Password is 5 Characters!"
	}

	if password1 != password2 {
		password2Err = "Password Confirmation is not match!"
	}

	if usernameErr == "" && emailErr == "" && password1Err == "" && password2Err == "" {
		if password1 == "" {
			newUser := models.User{
				Username: username,
				Email:    email,
			}
			err = models.DB.Model(&models.User{}).Where("user_id = ?", userId).Updates(&newUser).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/credential-information-page", c)
				return
			}
			middlewares.ClearSession(c)
			c.Redirect(
				http.StatusFound,
				"/shines/main/login-page",
			)
			return
		}
		newHashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(password1),
			bcrypt.DefaultCost,
		)
		if err != nil {
			panic(err)
		}
		newUser := models.User{
			Username: username,
			Email:    email,
			Password: string(newHashedPassword),
		}
		err = models.DB.Model(&models.User{}).Where("user_id = ?", userId).Updates(&newUser).Error
		if err != nil {

			ErrorHandler1("Failed to Update Data", "/shines/main/credential-information-page", c)
			return
		}
		middlewares.ClearSession(c)
		c.Redirect(
			http.StatusFound,
			"/shines/main/login-page",
		)
		return
	}
	context := gin.H{
		"title":        "Credential Information",
		"username":     middlewares.GetSession(c),
		"email":        user.Email,
		"image":        profile.Image,
		"password":     user.Password,
		"usernameErr":  usernameErr,
		"emailErr":     emailErr,
		"password1Err": password1Err,
		"password2Err": password2Err,
		"isSeller":     IsSeller(c),
		"isAdmin":      IsAdmin(c),
	}
	c.HTML(
		http.StatusOK,
		"profileCredential.html",
		context,
	)
}

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

func ViewCreateProductHandler(c *gin.Context) {
	role := GetRole(c)
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	if role == "Customer" {
		c.Redirect(
			http.StatusFound,
			"shines/main/home-page",
		)
		return
	}

	context := gin.H{
		"title":    "Create Product",
		"isSeller": IsSeller(c),
	}
	c.HTML(
		http.StatusOK,
		"createProduct.html",
		context,
	)
}

func CreateProductHandler(c *gin.Context) {
	var productNameErr, categoryErr, priceErr, quantityErr, fileErr string

	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}

	productName := c.PostForm("productName")
	description := c.PostForm("description")
	category := c.PostForm("category")
	priceString := c.PostForm("price")
	quantityString := c.PostForm("quantity")

	price, err := strconv.Atoi(priceString)
	if err != nil {
		priceErr = "Price must be a number!"
	}

	quantity, err := strconv.Atoi(quantityString)
	if err != nil {
		quantityErr = "Quantity must be a number!"
	}

	if price <= 0 {
		priceErr = "Price must be greater than 0!"
	}

	if quantity <= 0 {
		quantityErr = "Quantity must be greater than 0!"
	}

	if len(productName) < 3 {
		productNameErr = "Minimum Product Name is 3 Characters!"
	}

	if category == "" {
		categoryErr = "Category must be selected!"
	}

	file, err := c.FormFile("photo")
	if file == nil {
		if productNameErr == "" && categoryErr == "" && priceErr == "" && quantityErr == "" {
			sellerId := GetSellerId(c)
			product := models.Product{
				ShopId:             uint(sellerId),
				ProductName:        productName,
				ProductDescription: description,
				ProductCategory:    category,
				ProductPrice:       float64(price),
				ProductImage:       "productDefault.png",
				ProductStock:       uint(quantity),
			}
			err := models.DB.Create(&product).Error
			if err != nil {

				ErrorHandler1("Failed to Create Data", "/shines/main/create-product-page", c)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/shines/main/home-page",
			)
			return
		}
	} else {
		if err != nil {
			fileErr = "Failed Upload Picture"
		}
		err = c.SaveUploadedFile(file, "views/images/"+file.Filename)
		if err != nil {
			fileErr = "Failed Upload Picture"
		}

		if productNameErr == "" && categoryErr == "" && priceErr == "" && quantityErr == "" && fileErr == "" {
			sellerId := GetSellerId(c)
			product := models.Product{
				ShopId:             uint(sellerId),
				ProductName:        productName,
				ProductDescription: description,
				ProductCategory:    category,
				ProductPrice:       float64(price),
				ProductStock:       uint(quantity),
				ProductImage:       file.Filename,
			}
			err := models.DB.Create(&product).Error
			if err != nil {

				ErrorHandler1("Failed to Create Data", "/shines/main/create-product-page", c)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/shines/main/home-page",
			)
			return
		}
	}
	context := gin.H{
		"title":          "Create Product",
		"productNameErr": productNameErr,
		"categoryErr":    categoryErr,
		"priceErr":       priceErr,
		"quantityErr":    quantityErr,
		"fileErr":        fileErr,
		"productName":    productName,
		"description":    description,
		"category":       category,
		"price":          price,
		"isSeller":       IsSeller(c),
		"quantity":       quantity,
		"isAdmin":        IsAdmin(c),
	}
	c.HTML(
		http.StatusOK,
		"createProduct.html",
		context,
	)
}

func ViewUpdateProductHandler(c *gin.Context) {
	role := GetRole(c)
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	if role == "Customer" {
		c.Redirect(
			http.StatusFound,
			"shines/main/home-page",
		)
		return
	}

	productId := c.Param("productId")
	product := models.Product{}
	err := models.DB.Model(&models.Product{}).Select("*").Where("Product_id = ?", productId).First(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return
	}
	context := gin.H{
		"title":        "Update Product",
		"productName":  product.ProductName,
		"description":  product.ProductDescription,
		"category":     product.ProductCategory,
		"price":        product.ProductPrice,
		"productImage": product.ProductImage,
		"quantity":     product.ProductStock,
		"productId":    productId,
		"isSeller":     IsSeller(c),
		"isAdmin":      IsAdmin(c),
	}
	c.HTML(
		http.StatusOK,
		"updateProduct.html",
		context,
	)
}

func UpdateProductHandler(c *gin.Context) {
	var productNameErr, categoryErr, priceErr, quantityErr, fileErr string

	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}

	productName := c.PostForm("productName")
	description := c.PostForm("description")
	category := c.PostForm("category")
	priceString := c.PostForm("price")
	quantityString := c.PostForm("quantity")

	productId := c.Param("productId")
	product := models.Product{}
	err := models.DB.Model(&models.Product{}).Select("*").Where("Product_id = ?", productId).First(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return
	}

	price, err := strconv.Atoi(priceString)
	if err != nil {
		priceErr = "Price must be a number!"
	}

	quantity, err := strconv.Atoi(quantityString)
	if err != nil {
		quantityErr = "Quantity must be a number!"
	}

	if price <= 0 {
		priceErr = "Price must be greater than 0!"
	}

	if quantity <= 0 {
		quantityErr = "Quantity must be greater than 0!"
	}

	if len(productName) < 3 {
		productNameErr = "Minimum Product Name is 3 Characters!"
	}

	if category == "" {
		categoryErr = "Category must be selected!"
	}

	file, err := c.FormFile("photo")
	if file == nil {
		if productNameErr == "" && categoryErr == "" && priceErr == "" && quantityErr == "" {
			product.ProductName = productName
			product.ProductDescription = description
			product.ProductCategory = category
			product.ProductPrice = float64(price)
			product.ProductStock = uint(quantity)
			err = models.DB.Model(&models.Product{}).Where("Product_id = ?", productId).Updates(&product).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/home-page", c)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/shines/main/home-page",
			)
			return
		}
	} else {
		if err != nil {
			fileErr = "Failed Upload Picture"
		}
		err = c.SaveUploadedFile(file, "views/images/"+file.Filename)
		if err != nil {
			fileErr = "Failed Upload Picture"
		}

		if productNameErr == "" && categoryErr == "" && priceErr == "" && quantityErr == "" && fileErr == "" {
			product.ProductName = productName
			product.ProductDescription = description
			product.ProductCategory = category
			product.ProductPrice = float64(price)
			product.ProductStock = uint(quantity)
			product.ProductImage = file.Filename
			err = models.DB.Model(&models.Product{}).Where("Product_id = ?", productId).Updates(&product).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/home-page", c)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/shines/main/home-page",
			)
			return
		}
	}

	context := gin.H{
		"title":          "Update Product",
		"productNameErr": productNameErr,
		"categoryErr":    categoryErr,
		"priceErr":       priceErr,
		"quantityErr":    quantityErr,
		"fileErr":        fileErr,
		"productName":    productName,
		"description":    description,
		"category":       category,
		"price":          price,
		"quantity":       quantity,
		"isSeller":       IsSeller(c),
		"isAdmin":        IsAdmin(c),
		"productImage":   product.ProductImage,
	}
	c.HTML(
		http.StatusOK,
		"updateProduct.html",
		context,
	)
}

func ViewSellerCatalogHandler(c *gin.Context) {
	role := GetRole(c)
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	if role == "Customer" {
		c.Redirect(
			http.StatusFound,
			"shines/main/home-page",
		)
		return
	}

	shopId := GetShopId(c)
	products := []models.Product{}
	err := models.DB.Model(&models.Product{}).Select("*").Where("Shop_id = ?", shopId).Find(&products).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/seller-catalog-page", c)
		return
	}

	context := gin.H{
		"title":    "Seller Catalog",
		"products": products,
		"isSeller": IsSeller(c),
		"isAdmin":  IsAdmin(c),
	}
	c.HTML(
		http.StatusOK,
		"sellerCatalog.html",
		context,
	)
}

func ViewDeleteConfirmationHandler(c *gin.Context) {
	role := GetRole(c)
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	if role == "Customer" {
		c.Redirect(
			http.StatusFound,
			"shines/main/home-page",
		)
		return
	}

	productId := c.Param("productId")
	product := models.Product{}
	err := models.DB.Model(&models.Product{}).Select("*").Where("Product_id = ?", productId).First(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/seller-catalog-page", c)
		return
	}

	context := gin.H{
		"title":              "Delete Confirmation",
		"productName":        product.ProductName,
		"productId":          productId,
		"productPrice":       product.ProductPrice,
		"productStock":       product.ProductStock,
		"productImage":       product.ProductImage,
		"productDescription": product.ProductDescription,
		"isSeller":           IsSeller(c),
		"isAdmin":            IsAdmin(c),
	}
	c.HTML(
		http.StatusOK,
		"deleteConfirmation.html",
		context,
	)
}

func DeleteProductHandler(c *gin.Context) {
	role := GetRole(c)
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	if role == "Customer" {
		c.Redirect(
			http.StatusFound,
			"shines/main/home-page",
		)
		return
	}

	productId := c.Param("productId")
	product := models.Product{}
	err := models.DB.Model(&models.Product{}).Select("*").Where("Product_id = ?", productId).First(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/seller-catalog-page", c)
		return
	}
	DeleteProduct(c, productId)
	c.Redirect(
		http.StatusFound,
		"/shines/main/seller-catalog-page",
	)
}

func ViewAdminHandler(c *gin.Context) {
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
	users := []models.User{}
	err := models.DB.Model(&models.User{}).Select("*").Find(&users).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/administrator-page", c)
		return
	}
	context := gin.H{
		"title":    "Administrator",
		"users":    users,
		"isSeller": IsSeller(c),
		"isAdmin":  IsAdmin(c),
	}
	c.HTML(
		http.StatusOK,
		"admin.html",
		context,
	)
}

func ViewDetailPersonalHandler(c *gin.Context) {
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
	context := gin.H{
		"title":         "Detail User Information",
		"firstName":     profile.FirstName,
		"lastName":      profile.LastName,
		"address":       profile.Address,
		"image":         profile.Image,
		"userID":        userId,
		"isSeller":      IsSeller(c),
		"isAdmin":       IsAdmin(c),
		"isAdminTarget": IsAdminTarget(c, int(userId)),
	}
	c.HTML(
		http.StatusOK,
		"detailPersonal.html",
		context,
	)
}

func DetailPersonalHandler(c *gin.Context) {
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
	profile := models.Profile{}
	strUserId := c.Param("userId")
	UserId, _ := strconv.Atoi(strUserId)
	models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", UserId).First(&profile)
	var firstName, lastName, address string
	var firstNameErr, lastNameErr, addressErr, fileErr string

	firstName = c.PostForm("firstname")
	lastName = c.PostForm("lastname")
	address = c.PostForm("address")

	if len(firstName) < 2 && len(firstName) != 0 {
		firstNameErr = "Minimum First Name is 2 Character!"
	}
	if len(lastName) < 3 && len(lastName) != 0 {
		lastNameErr = "Minimum Last Name is 3 Characters!"
	}

	if len(address) < 5 && len(address) != 0 {
		addressErr = "Minimum Address is 5 Characters!"
	}

	file, err := c.FormFile("picture")
	if file == nil {
		if firstNameErr == "" && lastNameErr == "" && addressErr == "" {
			profile := models.Profile{
				UserID:    uint(UserId),
				FirstName: firstName,
				LastName:  lastName,
				Address:   address,
			}
			err := models.DB.Model(&models.Profile{}).Where("user_id = ?", UserId).Updates(&profile).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/personal-information-page", c)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/shines/main/administrator-page",
			)
			return
		}
		context := gin.H{
			"title":        "Detail User Information",
			"firstName":    profile.FirstName,
			"lastName":     profile.LastName,
			"address":      profile.Address,
			"firstNameErr": firstNameErr,
			"userID":       UserId,
			"lastNameErr":  lastNameErr,
			"image":        profile.Image,
			"addressErr":   addressErr,
		}
		c.HTML(
			http.StatusOK,
			"detailPersonal.html",
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
		if firstNameErr == "" && lastNameErr == "" && addressErr == "" && fileErr == "" {
			profile := models.Profile{
				UserID:    uint(UserId),
				FirstName: firstName,
				LastName:  lastName,
				Address:   address,
				Image:     file.Filename,
			}
			err := models.DB.Model(&models.Profile{}).Where("user_id = ?", UserId).Updates(&profile).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/personal-information-page", c)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/shines/main/administrator-page",
			)
			return
		}
		context := gin.H{
			"title":         "Detail User Information",
			"firstName":     profile.FirstName,
			"lastName":      profile.LastName,
			"address":       profile.Address,
			"image":         profile.Image,
			"firstNameErr":  firstNameErr,
			"lastNameErr":   lastNameErr,
			"addressErr":    addressErr,
			"fileErr":       fileErr,
			"isSeller":      IsSeller(c),
			"isAdmin":       IsAdmin(c),
			"isAdminTarget": IsAdminTarget(c, int(UserId)),
		}
		c.HTML(
			http.StatusOK,
			"detailPersonal.html",
			context,
		)
	}
}

func ViewDetailCredentialHandler(c *gin.Context) {
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
	user := models.User{}
	err := models.DB.Model(&models.User{}).Select("*").Where("User_id = ?", userId).First(&user).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/administrator-page", c)
		return
	}
	profile := models.Profile{}
	err = models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", userId).First(&profile).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/administrator-page", c)
		return
	}
	fmt.Println(IsAdminTarget(c, userId))
	context := gin.H{
		"title":         "Detail Credential Information",
		"username":      user.Username,
		"email":         user.Email,
		"image":         profile.Image,
		"firstName":     profile.FirstName,
		"lastName":      profile.LastName,
		"address":       profile.Address,
		"isSeller":      IsSeller(c),
		"role":          user.Role,
		"userID":        userId,
		"isAdmin":       IsAdmin(c),
		"isAdminTarget": IsAdminTarget(c, userId),
	}
	c.HTML(
		http.StatusOK,
		"detailCredential.html",
		context,
	)
}

func DetailCredentialHandler(c *gin.Context) {
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
	strUserID := c.Param("userId")
	userId, _ := strconv.Atoi(strUserID)

	user := models.User{}
	err := models.DB.Model(&models.User{}).Select("*").Where("User_id = ?", userId).First(&user).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/credential-information-page", c)
		return
	}
	profile := models.Profile{}
	err = models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", userId).First(&profile).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/credential-information-page", c)
		return
	}
	var username, email, password1, password2 string
	var usernameErr, emailErr, password1Err, password2Err string

	username = c.PostForm("username")
	email = c.PostForm("email")
	roleTarget := c.PostForm("role")
	password1 = c.PostForm("password")
	password2 = c.PostForm("password-confirmation")

	if len(username) < 5 {
		usernameErr = "Minimum Username is 5 Characters!"
	}
	if len(email) < 5 {
		emailErr = "Minimum Email is 5 Characters!"
	}
	if password1 != "" && len(password1) < 5 {
		password1Err = "Minimum Password is 5 Characters!"
	}

	if password1 != password2 {
		password2Err = "Password Confirmation is not match!"
	}

	if usernameErr == "" && emailErr == "" && password1Err == "" && password2Err == "" {
		if password1 == "" {
			newUser := models.User{
				Username: username,
				Email:    email,
				Role:     models.UserRole(roleTarget),
			}
			err = models.DB.Model(&models.User{}).Where("user_id = ?", userId).Updates(&newUser).Error
			if err != nil {

				ErrorHandler1("Failed to Update Data", "/shines/main/credential-information-page", c)
				return
			}
			targetUrl := fmt.Sprintf("/shines/main/detail-credential-page/%d", userId)
			c.Redirect(
				http.StatusFound,
				targetUrl,
			)
			return
		}
		newHashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(password1),
			bcrypt.DefaultCost,
		)
		if err != nil {
			panic(err)
		}
		newUser := models.User{
			Username: username,
			Email:    email,
			Role:     models.UserRole(roleTarget),
			Password: string(newHashedPassword),
		}
		err = models.DB.Model(&models.User{}).Where("user_id = ?", userId).Updates(&newUser).Error
		if err != nil {

			ErrorHandler1("Failed to Update Data", "/shines/main/credential-information-page", c)
			return
		}
		targetUrl := fmt.Sprintf("/shines/main/detail-credential-page/%d", userId)
		c.Redirect(
			http.StatusFound,
			targetUrl,
		)
		return
	}
	context := gin.H{
		"title":         "Credential Information",
		"username":      middlewares.GetSession(c),
		"email":         user.Email,
		"image":         profile.Image,
		"password":      user.Password,
		"usernameErr":   usernameErr,
		"emailErr":      emailErr,
		"password1Err":  password1Err,
		"password2Err":  password2Err,
		"isSeller":      IsSeller(c),
		"isAdmin":       IsAdmin(c),
		"isAdminTarget": IsAdminTarget(c, userId),
	}
	c.HTML(
		http.StatusOK,
		"detailCredential.html",
		context,
	)
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

func ViewDetailProductHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	productId := c.Param("productId")
	product := models.Product{}
	err := models.DB.Model(&models.Product{}).Select("*").Where("product_id = ?", productId).First(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return
	}
	shopId := product.ShopId
	shop := models.Shop{}
	err = models.DB.Model(&models.Shop{}).Select("*").Where("seller_id = ?", shopId).First(&shop).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return
	}
	stockSlice := make([]int, product.ProductStock)
	for i := 0; i < int(product.ProductStock); i++ {
		stockSlice[i] = i + 1
	}
	context := gin.H{
		"title":         "Detail Product",
		"productName":   product.ProductName,
		"description":   product.ProductDescription,
		"category":      product.ProductCategory,
		"price":         product.ProductPrice,
		"productImage":  product.ProductImage,
		"stock":         product.ProductStock,
		"shopId":        shopId,
		"shopName":      shop.ShopName,
		"quantityOrder": stockSlice,
		"isSeller":      IsSeller(c),
		"isAdmin":       IsAdmin(c),
	}
	c.HTML(
		http.StatusOK,
		"detailProduct.html",
		context,
	)
}

func DetailProductHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	strProductId := c.Param("productId")
	productId, _ := strconv.Atoi(strProductId)
	strOrderQuantity := c.PostForm("quantity")
	orderQuantity, _ := strconv.Atoi(strOrderQuantity)

	product := models.Product{}
	err := models.DB.Model(&models.Product{}).Select("*").Where("product_id = ?", productId).First(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return
	}

	AddToCart(c, productId, orderQuantity, int(product.ProductStock))
	c.Redirect(
		http.StatusFound,
		"/shines/main/home-page",
	)
}

func ViewCartHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	userId := GetuserId(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	cart := []models.Cart{}
	err := models.DB.Model(&models.Cart{}).Select("*").Where("user_id = ?", userId).Find(&cart).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return

	}
	totalPrice := new(big.Float).SetFloat64(0.0)
	for _, item := range cart {
		priceProduct := new(big.Float).SetFloat64(GetPriceProduct(c, int(item.ProductID)))
		quantity := big.NewFloat(float64(item.Quantity))
		totalPrice.Add(totalPrice, new(big.Float).Mul(priceProduct, quantity))

	}
	transactions := []models.TransactionDetail{}
	err = models.DB.Table("carts").
		Select("carts.cart_id, carts.user_id, users.username, users.email, carts.product_id, products.product_name as product_name, products.product_price as price, carts.quantity").
		Joins("left join users on carts.user_id = users.user_id").
		Joins("left join products on carts.product_id = products.product_id").
		Where("carts.user_id = ?", userId).
		Find(&transactions).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return

	}
	context := gin.H{
		"title":        "Cart",
		"totalPrice":   totalPrice,
		"isSeller":     IsSeller(c),
		"Transactions": transactions,
		"isAdmin":      IsAdmin(c),
	}
	c.HTML(
		http.StatusOK,
		"cart.html",
		context,
	)
}

func ViewUpdateCartHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	strCartId := c.Param("cartId")
	cartId, _ := strconv.Atoi(strCartId)
	cart := models.Cart{}
	err := models.DB.Model(&models.Cart{}).Select("*").Where("cart_id = ?", cartId).First(&cart).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/cart-page", c)
		return
	}
	productId := int(cart.ProductID)
	product := models.Product{}
	err = models.DB.Model(&models.Product{}).Select("*").Where("product_id = ?", productId).First(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/cart-page", c)
		return
	}
	stockSlice := make([]int, product.ProductStock)
	for i := 0; i < int(product.ProductStock); i++ {
		stockSlice[i] = i + 1
	}
	productId = int(cart.ProductID)
	shop := models.Shop{}
	err = models.DB.Model(&models.Shop{}).Select("*").Where("seller_id = ?", product.ShopId).First(&shop).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/cart-page", c)

		return

	}
	context := gin.H{
		"title":         "Update Cart",
		"productName":   product.ProductName,
		"description":   product.ProductDescription,
		"category":      product.ProductCategory,
		"price":         product.ProductPrice,
		"shopName":      shop.ShopName,
		"productPrice":  product.ProductPrice,
		"productImage":  product.ProductImage,
		"quantity":      cart.Quantity,
		"stock":         product.ProductStock,
		"cartId":        cartId,
		"quantityOrder": stockSlice,
		"isSeller":      IsSeller(c),
		"isAdmin":       IsAdmin(c),
	}
	c.HTML(
		http.StatusOK,
		"updateCart.html",
		context,
	)
}

func UpdateCartHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	strCartId := c.Param("cartId")
	cartId, _ := strconv.Atoi(strCartId)
	strOrderQuantity := c.PostForm("quantity")
	orderQuantity, _ := strconv.Atoi(strOrderQuantity)
	cart := models.Cart{}
	err := models.DB.Model(&models.Cart{}).Select("*").Where("cart_id = ?", cartId).First(&cart).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/cart-page", c)
		return
	}
	productId := int(cart.ProductID)
	product := models.Product{}
	err = models.DB.Model(&models.Product{}).Select("*").Where("product_id = ?", productId).First(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/cart-page", c)
		return
	}
	fmt.Println(orderQuantity)
	stock := int(product.ProductStock)
	UpdateCart(c, cartId, orderQuantity, stock)
	c.Redirect(
		http.StatusFound,
		"/shines/main/cart-page",
	)
}

func DeleteCartHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	strCartId := c.Param("cartId")
	cartId, _ := strconv.Atoi(strCartId)
	DeleteCart(c, cartId)
	c.Redirect(
		http.StatusFound,
		"/shines/main/cart-page",
	)
}

func CheckoutHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
		http.StatusFound,
		"shines/main/login-page",
		)
		return
	}
	userId := GetuserId(c)
	cart := []models.Cart{}
	err := models.DB.Model(&models.Cart{}).Select("*").Where("user_id = ?", userId).Find(&cart).Error
	if err != nil {
		context := gin.H {
			"title":"Error",
			"message":"Failed to Get Data",
			"source":"/shines/main/cart-page",
		}
		c.HTML(
			http.StatusInternalServerError,
			"error.html",
			context,
		)
		return
	}
	details := []models.TransactionDetail{}
	err = models.DB.Table("carts").
	Select("carts.cart_id, carts.user_id, users.username, users.email, carts.product_id, products.product_name as product_name, products.product_price as price, carts.quantity as quantity").
	Joins("left join users on carts.user_id = users.user_id").
	Joins("left join products on carts.product_id = products.product_id").
	Where("carts.user_id = ?", userId).
	Find(&details).Error
	if err != nil {
		context := gin.H {
			"title":"Error",
			"message":"Failed to Get Data",
			"source":"/shines/main/cart-page",
		}
		c.HTML(
			http.StatusInternalServerError,
			"error.html",
			context,
		)
		return
	}
	for _, item := range details {
		AddToTransaction(c, item.Price, int(item.ProductID), int(item.Quantity))
		UpdateStockProduct(c, int(item.ProductID), int(item.Quantity))
	}
	ClearCart(c)
	c.Redirect(
		http.StatusFound,
		"/shines/main/home-page",
	)
}