package controllers

import (
	"math/big"
	"net/http"
	"shines/middlewares"
	"shines/models"
	"shines/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ViewCartHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	buyerId := repositories.GetuserId(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	cart := []models.Cart{}
	err := models.DB.Model(&models.Cart{}).Select("*").Where("buyer_id = ?", buyerId).Find(&cart).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return

	}
	totalPrice := new(big.Float).SetFloat64(0.0)
	for _, item := range cart {
		priceProduct := new(big.Float).SetFloat64(repositories.GetPriceProduct(c, int(item.ProductID)))
		quantity := big.NewFloat(float64(item.Quantity))
		totalPrice.Add(totalPrice, new(big.Float).Mul(priceProduct, quantity))

	}
	transactions := []models.TransactionDetail{}
	err = models.DB.Table("carts").
		Select("carts.cart_id, carts.buyer_id, users.username, users.email, carts.product_id, products.product_name as product_name, products.product_price as price, carts.quantity").
		Joins("left join users on carts.buyer_id = users.user_id").
		Joins("left join products on carts.product_id = products.product_id").
		Where("carts.buyer_id = ?", buyerId).
		Find(&transactions).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return

	}
	context := gin.H{
		"title":        "Cart",
		"totalPrice":   totalPrice,
		"isSeller":     repositories.IsSeller(c),
		"Transactions": transactions,
		"isAdmin":      repositories.IsAdmin(c),
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
		"isSeller":      repositories.IsSeller(c),
		"isAdmin":       repositories.IsAdmin(c),
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
	stock := int(product.ProductStock)
	var urlSource string
	err, urlSource = repositories.UpdateCart(c, cartId, orderQuantity, stock)
	if err != nil {

		ErrorHandler1("Failed to Get Data", urlSource, c)

	}
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

	var urlSource string
	err, urlSource := repositories.DeleteCart(c, cartId)
	if err != nil {

		ErrorHandler1("Failed to Delete Data", urlSource, c)
	}
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
	buyerId := repositories.GetuserId(c)
	cart := []models.Cart{}
	err := models.DB.Model(&models.Cart{}).Select("*").Where("buyer_id = ?", buyerId).Find(&cart).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/cart-page", c)
		return
	}
	details := []models.TransactionDetail{}
	err = models.DB.Table("carts").
		Select("carts.cart_id, carts.buyer_id, users.username, users.email, carts.product_id, products.product_name as product_name, products.product_price as price, carts.quantity as quantity").
		Joins("left join users on carts.buyer_id = users.user_id").
		Joins("left join products on carts.product_id = products.product_id").
		Where("carts.buyer_id = ?", buyerId).
		Find(&details).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/cart-page", c)
		return
	}
	for _, item := range details {
		err := repositories.AddToTransaction(c, item.Price, int(item.ProductID), int(item.Quantity))
		if err != nil {
			ErrorHandler1("Failed to Create Data", "/shines/main/cart-page", c)
			return
		}

		err = repositories.UpdateStockProduct(c, int(item.ProductID), int(item.Quantity))
		if err != nil {

			ErrorHandler1("Failed to Update Data", "/shines/main/cart-page", c)
			return
		}
	}
	err = repositories.ClearCart(c)
	if err != nil {

		ErrorHandler1("Failed to Delete Data", "/shines/main/cart-page", c)
		return
	}
	c.Redirect(
		http.StatusFound,
		"/shines/main/home-page",
	)
}