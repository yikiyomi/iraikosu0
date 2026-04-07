package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func  () {
  router := gin.Default()
  router.GET("/user/:name", func(c *gin.Context) {
    name := c.Param("name")
    c.String(http.StatusOK, "hello %s喵", name)
  })
  router.GET("/user/:name/*action", func(c *gin.Context) {
    name := c.Param("name")
    action := c.Param("action")
    message := name + "是" + action
    c.String(http.StatusOK, message)
  })
  router.Run(":8080")
}