package controllers

import (
	"fmt"
	"net/http"
	"shines/middlewares"
	"shines/models"

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
	err := models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", GetIdUser(c)).First(&profile).Error
	if err != nil {
		context := gin.H {
			"title":"Error",
			"message":"Failed to Get Data",
			"source":"/shines/main/credential-information-page",
		}
		c.HTML(
			http.StatusInternalServerError,
			"error.html",
			context,
		)
		return
	}
	username := middlewares.GetSession(c)
	context := gin.H {
		"title":"Credential Information",
		"username":username,
		"image":profile.Image,
		"firstName":profile.FirstName,
		"lastName":profile.LastName,
		"address":profile.Address,
		"email":email,
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
	err := models.DB.Model(&models.User{}).Select("*").Where("User_id = ?", GetIdUser(c)).First(&user).Error
	if err != nil {
		context := gin.H {
			"title":"Error",
			"message":"Failed to Get Data",
			"source":"/shines/main/credential-information-page",
		}
		c.HTML(
			http.StatusInternalServerError,
			"error.html",
			context,
		)
		return
	}
	profile := models.Profile{}
	err = models.DB.Model(&models.Profile{}).Select("*").Where("User_id = ?", GetIdUser(c)).First(&profile).Error
	if err != nil {
		context := gin.H {
			"title":"Error",
			"message":"Failed to Get Data",
			"source":"/shines/main/credential-information-page",
		}
		c.HTML(
			http.StatusInternalServerError,
			"error.html",
			context,
		)
		return
	}
	var username, email, password1, password2 string
	var usernameErr, emailErr, password1Err, password2Err string

	userId := GetIdUser(c)

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
			newUser := models.User {
				Username: username,
				Email: email,
			}
			err = models.DB.Model(&models.User{}).Where("user_id = ?", userId).Updates(&newUser).Error
			if err != nil {
				context := gin.H{
					"title":   "Error",
					"message": "Failed to Update Data",
					"source":  "/shines/main/credential-information-page",
				}
				c.HTML(
					http.StatusInternalServerError,
					"error.html",
					context,
				)
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
		newUser := models.User {
			Username: username,
			Email: email,
			Password: string(newHashedPassword),
		}
		err = models.DB.Model(&models.User{}).Where("user_id = ?", userId).Updates(&newUser).Error
		if err != nil {
			context := gin.H{
				"title":   "Error",
				"message": "Failed to Update Data",
				"source":  "/shines/main/credential-information-page",
			}
			c.HTML(
				http.StatusInternalServerError,
				"error.html",
				context,
			)
			return
		}
		middlewares.ClearSession(c)
		c.Redirect(
			http.StatusFound,
			"/shines/main/login-page",
		)
		return
	}
	context := gin.H {
		"title":"Credential Information",
		"username":middlewares.GetSession(c),
		"email":user.Email,
		"image":profile.Image,
		"password":user.Password,
		"usernameErr":usernameErr,
		"emailErr":emailErr,
		"password1Err":password1Err,
		"password2Err":password2Err,
	}
	c.HTML(
		http.StatusOK,
		"profileCredential.html",
		context,
	)
}