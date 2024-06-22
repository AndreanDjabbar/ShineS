package controllers

import (
	"fmt"
	"math/big"
	"net/http"
	"shines/middlewares"
	"shines/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

		ErrorHandler1("Failed to Get Data", "/shines/main/cart-page", c)
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

		ErrorHandler1("Failed to Get Data", "/shines/main/cart-page", c)
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