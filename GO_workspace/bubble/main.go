package main

import (
	"net/http"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
)
//  Todo model
type Todo struct{
	ID int `json:"id"`
	Title string `json:""tltle`
	Status bool `json:"status"`
}

func initmysql()(err error){
	dsn :="root:123456@tcp(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	DB,err = gorm.Open("mysql", dsn)
	if err!= nil {
		return 
	}
	return DB.DB().Ping()
}

func main(){
	// 创建数据库
	// 在dg中创建bubble数据库
	//连接数据库 
	err:= initmysql()
	if err!= nil {
		panic(err)
	}
	defer DB.Close()
	// 模型绑定
	DB.AutoMigrate(&Todo{})
	r:=gin.Default()
	// 告诉gin框架模板文件引用的静态文件夹
	r.Static("/static","static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/",func(c *gin.Context) {
		c.HTML(http.StatusOK,"index.html",nil)
	})
	v1Group :=r.Group("/v1")
	{
	//待办事项

	//添加
	v1Group.POST("/todo",func(c *gin.Context) {

	//1..拿去数据
	var todo Todo
	c.BindJSON(&todo)

	// 2..存入数据库并返回响应
		if err :=DB.Create(&todo).Error;err != nil{
			c.JSON(http.StatusOK,gin.H{"error":err.Error()})
		}else{
			c.JSON(http.StatusOK,todo)
		}

	})

	//查看所有代办事项
	v1Group.GET("/todo",func(c *gin.Context) {
		//查询所有事项
		var todolist []Todo
		if err =DB.Find(&todolist).Error;err!=nil{
			c.JSON(http.StatusOK,gin.H{"error":err.Error()})
		}else{
			c.JSON(http.StatusOK,todolist)
		}
	})
	//查看一个代办事项
	v1Group.GET("/todo/:id",func(c *gin.Context) {

	})
	//修改
	v1Group.PUT("/todo/:id",func(c *gin.Context) {
		id, ok :=c.Params.Get("id")
		if !ok{
			c.JSON(http.StatusOK,gin.H{"error":"无效id喵"})
			return 
		}
		var todo Todo
		if err=DB.Where("id=?",id).First(&todo).Error;err!= nil{
			c.JSON(http.StatusOK,gin.H{"error":err.Error()})
			return
		}
		c.BindJSON(&todo)
		if err=DB.Save(&todo).Error;err!= nil{
			c.JSON(http.StatusOK,gin.H{"errror":err.Error()})
		}else{
			c.JSON(http.StatusOK,todo)
		}
	})
	//删除
	v1Group.DELETE("/todo/:id",func(c *gin.Context) {
		id, ok :=c.Params.Get("id")
		if !ok{
			c.JSON(http.StatusOK,gin.H{"error":"无效id喵"})
			return 
		}
		if err =DB.Where("id=?",id).Delete(Todo{}).Error;err!= nil{
			c.JSON(http.StatusOK,gin.H{"error":err.Error()})
		}else{
			c.JSON(http.StatusOK,gin.H{"message":"删除成功喵"})
		}
	})
	}
	r.Run()
}