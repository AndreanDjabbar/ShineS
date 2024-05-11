package controllers

import (
	"fmt"
	"net/http"
	"shines/middlewares"
	"shines/models"

	"github.com/gin-gonic/gin"
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
	user := middlewares.GetSession(c)
	isLogged := middlewares.CheckSession(c)
	fmt.Println(isLogged)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	context := gin.H {
		"title":"Home",
		"user":user,
	}
	c.HTML(
		http.StatusOK,
		"home.html",
		context,
	)
}

func GetIdUser(c *gin.Context) int {
	user := middlewares.GetSession(c)
	var userId int
	models.DB.Model(&models.User{}).Select("UserId").Where("username = ?", user).First(&userId)
	return userId
}

func CreateProfile(c *gin.Context) {
	userId := GetIdUser(c)

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
	models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", GetIdUser(c)).First(&profile)
	context := gin.H {
		"title":"Personal Information",
		"image":profile.Image,
		"firstName":profile.FirstName,
		"lastName":profile.LastName,
		"address":profile.Address,
	}
	c.HTML(
		http.StatusOK,
		"profilePersonal.html",
		context,
	)
}

func PersonalHandler(c *gin.Context) {
	profile := models.Profile{}
	models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", GetIdUser(c)).First(&profile)
	var firstName, lastName, address string
	var firstNameErr, lastNameErr, addressErr, fileErr string
	userId := GetIdUser(c)

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
		if firstNameErr == "" && lastNameErr == ""  && addressErr == "" {
			profile := models.Profile {
				UserID: uint(userId),
				FirstName: firstName,
				LastName: lastName,
				Address: address,
			}
			err := models.DB.Model(&models.Profile{}).Where("user_id = ?", userId).Updates(&profile).Error
			if err != nil {
				context := gin.H{
					"title":   "Error",
					"message": "Failed to Update Data",
					"source":  "/shines/main/personal-information-page",
				}
				c.HTML(
					http.StatusInternalServerError,
					"error.html",
					context,
				)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/shines/main/personal-information-page",
			)
			return
		}
		context := gin.H {
			"title":"Personal Information",
			"firstName":profile.FirstName,
			"lastName":profile.LastName,
			"address":profile.Address,
			"firstNameErr":firstNameErr,
			"lastNameErr":lastNameErr,
			"image":profile.Image,
			"addressErr":addressErr,
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
		if firstNameErr == "" && lastNameErr == ""  && addressErr == "" && fileErr == "" {
			profile := models.Profile {
				UserID: uint(userId),
				FirstName: firstName,
				LastName: lastName,
				Address: address,
				Image: file.Filename,
			}
			err := models.DB.Model(&models.Profile{}).Where("user_id = ?", userId).Updates(&profile).Error
			if err != nil {
				context := gin.H{
					"title":   "Error",
					"message": "Failed to Update Data",
					"source":  "/shines/main/personal-information-page",
				}
				c.HTML(
					http.StatusInternalServerError,
					"error.html",
					context,
				)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/shines/main/personal-information-page",
			)
			return
		}
		context := gin.H {
			"title":"Personal Information",
			"firstName":profile.FirstName,
			"lastName":profile.LastName,
			"address":profile.Address,
			"image":profile.Image,
			"firstNameErr":firstNameErr,
			"lastNameErr":lastNameErr,
			"addressErr":addressErr,
			"fileErr":fileErr,
		}
		c.HTML(
			http.StatusOK,
			"profilePersonal.html",
			context,
		)
	}
}