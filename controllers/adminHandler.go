package controllers

import (
	"net/http"
	"shines/middlewares"
	"shines/models"

	"github.com/gin-gonic/gin"
)

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