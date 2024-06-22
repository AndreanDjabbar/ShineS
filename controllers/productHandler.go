package controllers

import (
	"net/http"
	"shines/middlewares"
	"shines/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	err = DeleteProduct(c, productId)
	if err != nil {

		ErrorHandler1("Failed to Delete Data", "/shines/main/seller-catalog-page", c)
		return
	}
	c.Redirect(
		http.StatusFound,
		"/shines/main/seller-catalog-page",
	)
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

	sellerId := GetSellerIdByProductId(c, productId)
	product := models.Product{}
	err := models.DB.Model(&models.Product{}).Select("*").Where("product_id = ?", productId).First(&product).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return
	}

	var urlSource string
	err, urlSource = AddToCart(c, sellerId, productId, orderQuantity, int(product.ProductStock))
	if err != nil {

		ErrorHandler1("Failed to Create Data", urlSource, c)
		return
	}

	c.Redirect(
		http.StatusFound,
		"/shines/main/home-page",
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

	err, shopId := GetShopId(c)
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/shop-information-page", c)
		return
	}
	products := []models.Product{}
	err = models.DB.Model(&models.Product{}).Select("*").Where("Shop_id = ?", shopId).Find(&products).Error
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
