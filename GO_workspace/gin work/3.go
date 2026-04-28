package main
import(
	"net/http"
	"github.com/gin-gonic/gin"
)
type Response struct{
	Code int `json:"code"`
	Msg string`json:"msg"`
	Data interface{} `json:"data"`
}
type User struct{
	ID uint `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
}
var userList =[]User{
	{ID:1,Name:"张三",Age:20},
	{ID:2,Name:"李四",Age:18},
}
func main(){
	r:=gin.Default()

	// 路由组
	v1:=r.Group("miaomiao/v1")
	{
		// restdul路由
		v1.GET("/user",GetUserList)
		// 获取所有用户
		v1.GET("/user/:id",GetUserbyid)
		// 获取单个用户
		v1.POST("/user",CreateUser)
		// 新增用户
		v1.PUT("/user/:id",UpdataUser)
		// 全量更新
		v1.DELETE(".user/:id",Deleteuser)
		// 删除用户
	}
	r.Run(":8080")
}
// /获取用户列表
func GetUserList(c *gin.Context){
	c.JSON(http.StatusOK,Response{
		Msg:"查询成功",
		Data: userList,
	})
}

// 获取单个用户
func GetUserbyid(c *gin.Context){
	id:=c.Param("id")
	c.JSON(http.StatusOK,Response{
		Msg:"查询成功",
		Data: id,
	})
}
// 创建用户
func CreateUser(c *gin.Context){
	var user User
	if err:=c.ShouldBind(&user);err!=nil{
		c.JSON(http.StatusBadRequest,Response{
			Code:400,
			Msg: "参数错误",
			Data: nil,
		})
		return 
	}
	userList=append(userList, user)
	c.JSON(http.StatusCreated,Response{
		Code: 201,
		Msg: "创建成功",
		Data: user,
	})
}
// 更新用户
func UpdataUser(c *gin.Context){
	id:=c.Param("id")
	c.JSON(http.StatusOK,Response{
		Code: 200,
		Msg: "更新成功",
		Data: id ,
	})
}
// 删除用户
func Deleteuser(c *gin.Context){
	id:=c.Param("id")
	c.JSON(http.StatusOK,Response{
		Code: 200,
		Msg: "成功",
		Data: id,
	})
}