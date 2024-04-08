package controllers

import (
	"fmt"
	_ "fmt"
	"net/http"
	"shines/middlewares"
	"shines/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ViewRegisterHandler(c *gin.Context) {
		fmt.Println(middlewares.CheckSession(c))
		fmt.Println(middlewares.GetSession(c))
	fmt.Println(middlewares.CheckSession(c))
	context := gin.H {
		"title":"Sign Up",
	}
	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}

func isNumber (strings string) bool {
	for a := 0; a < len(strings); a++ {
		_, err := strconv.Atoi(string(strings[a]))
		if err != nil {
			return false
		}
	}
	return true
}


func RegisterHandler(c *gin.Context) {
	var user models.User

	var username, email, phone, password string
	var usernameErr, phoneErr, passwordErr, emailErr string

	username = c.PostForm("name")
	email = c.PostForm("email")
	phone = c.PostForm("phone")
	password = c.PostForm("password")

	if len(username) < 5 {
		usernameErr = "Minimum Username is 5 Characters!"
	}

	if !strings.Contains(email, "@") {
		emailErr = "Email must included @"
	}	

	if len(email) < 10 {
		emailErr = "Email must be at least 10 Characters and included @"
	}

	if len(phone) < 8 {
		phoneErr = "Minimum phone is 8 Characters!"
	}

	if !isNumber(phone) {
		phoneErr = "Phone must be a number"
	}

	if len(password) < 5 {
		passwordErr = "Minimum Password is 5 Characters!"
	}

	if usernameErr == "" && phoneErr == "" && passwordErr == "" && emailErr == "" {
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			panic(err)
		}

		user = models.User{
			Username: username,
			Email:    email,
			Phone:    phone,
			Password: string(hashedPassword),
			Role: "Customer",
		}

		err = models.DB.Create(&user).Error
		if err != nil {
			context := gin.H{
				"title":   "Error Create",
				"message": "Failed to Create Data",
				"source":  "/shines/main/register",
			}
			c.HTML(
				http.StatusOK,
				"error.html",
				context,
			)
		}
		c.Redirect(
			http.StatusTemporaryRedirect,
			"/shines/main/login",
		)
}

	context := gin.H{
		"title":"Sign Up",
		"email":email,
		"phone":phone,
		"username":username,
		"usernameErr": usernameErr,
		"emailErr":    emailErr,
		"phoneErr":    phoneErr,
		"passwordErr": passwordErr,
	}

	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}

func ViewLoginHandler(c *gin.Context) {
	
	fmt.Println(middlewares.CheckSession(c))
	fmt.Println(middlewares.GetSession(c))
	context := gin.H {
		"title":"Login",
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func LoginHandler(c *gin.Context) {
	var user models.User
	var username, password string
	var usernameErr, passwordErr string

	username = c.PostForm("username")
	password = c.PostForm("password")

	if len(username) < 5 {
		usernameErr = "Minimum Username is 5 Characters!"
	}

	if len(password) < 5 {
		passwordErr = "Minimum Password is 5 Characters!"
	}

	err := models.DB.Where("Username = ?", username).First(&user).Error
	if err != nil {
		usernameErr = "Invalid Username"
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		passwordErr = "Invalid Password"
	}

	if usernameErr == "" && passwordErr == "" {
		middlewares.SaveSession(c, username)
		c.Redirect(
			http.StatusFound,
			"/shines/main/homes",
		)
		return
	}

	context := gin.H {
		"title":"Login",
		"username":username,
		"usernameErr":usernameErr,
		"passwordErr":passwordErr,
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func LogoutHandler(c *gin.Context) {
	middlewares.ClearSession(c)
	c.Redirect(
		http.StatusFound,
		"/shines/main/login",
	)
}