package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handling all error message

func ErrorHandler1(message string, source string, c *gin.Context) {
	context := gin.H{
		"title":   "Error",
		"message": message,
		"source":  source,
	}
	c.HTML(
		http.StatusInternalServerError,
		"error.html",
		context,
	)
}

func ErrorHandler2(title string, message string, source string, c *gin.Context) {
	context := gin.H{
		"title":   title,
		"message": message,
		"source":  source,
	}
	c.HTML(
		http.StatusOK,
		"error.html",
		context,
	)
}
