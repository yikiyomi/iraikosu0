package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getting(r *gin.Context) {
	r.JSON(http.StatusOK, gin.H{"miaomiao": "GET"})
}

func posting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"miaomiao": "POST"})
}

func putting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"miaomiao": "PUT"})
}
func deleteing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"miaomiao": "删除成功喵DELETE"})
}
func patching(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"miaomiao": "PATCH"})
}
func head(c *gin.Context) {
	c.Status(http.StatusOK)
}

func options(c *gin.Context) {
	c.Status(http.StatusOK)
}
func mai
() {

	r := gin.Default()

	r.GET("/someGet", getting)
	r.POST("/somePost", posting)
	r.PUT("/somePut", putting)
	r.DELETE("/someDelete", deleteing)
	r.PATCH("/somePatch", patching)
	r.HEAD("/someHead", head)
	r.OPTIONS("/someOptions", options)
	r.Run()
}

//get  post  put  patch  delete  head  options
