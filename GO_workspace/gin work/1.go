package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func () {
  router := gin.Default()

//   router.GET("/welcome", func(c *gin.Context) {
//     firstname := c.DefaultQuery("firstname", "喵喵喵")
//     lastname := c.Query("lastname") 
//     c.String(http.StatusOK, "Hello %s %s", lastname,firstname )
//   })

 router.GET("/welcome",func(c *gin.Context){
	fname:=c.DefaultQuery("firstname","喵喵喵")
	lname:=c.Query("lastname")
	c.String(http.StatusOK,"hello %s %s",lname,fname)
  })

  router.Run(":8080")
}