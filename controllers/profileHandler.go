package controllers

import (
	"fmt"
	"net/http"
	"shines/middlewares"
	"shines/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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

func ViewHistoryHandler(c *gin.Context) {
	transactions := []models.Transactions{}
	salesTransactions := []models.Transactions{}
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	userId := GetuserId(c)
	role := GetRole(c)
	if role == "Admin" {
		transactions = []models.Transactions{}
		err := models.DB.Model(&models.Transactions{}).Select("*").Find(&transactions).Error
		if err != nil {
			ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
			return
		}
	} else {
		transactions = []models.Transactions{}
		err := models.DB.Model(&models.Transactions{}).Select("*").Where("Buyer_id = ?", userId).Find(&transactions).Error
		if err != nil {
			ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
			return
		}
		salesTransactions = []models.Transactions{}
		err = models.DB.Model(&models.Transactions{}).Select("*").Where("Seller_id = ?", userId).Find(&salesTransactions).Error
		if err != nil {
			ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
			return
		}
	}
	context := gin.H{
		"title":  "History",
		"userId": userId,
		"role":   role,
		"isSeller": IsSeller(c),
		"isAdmin":  IsAdmin(c),
		"transactions": transactions,
		"salesTransactions": salesTransactions,
	}
	c.HTML(
		http.StatusOK,
		"history.html",
		context,
	)
}