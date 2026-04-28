package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//全局变量
var (
	db *gorm.DB
	rdb *redis.Client
	ctx = context.Background()
)
// jwt密钥
var jwtKey =[]byte("my_secret_key")

func main() {
// 连接mysql
	dsn := "root:123456@tcp(127.0.0.1:3306)/allto?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db,err =gorm.Open(mysql.Open(dsn),&gorm.Config{})	
	if err !=nil{
		log.Fatal("mysql连接失败:",err)
	}
// 连接redis
rdb =redis.NewClient(&redis.Options{
	Addr:  "localhost:6379",
	Password: "",
	DB: 0,
})
if err := rdb.Ping(ctx).Err(); err != nil {
	log.Fatal("redis连接失败:", err)	
}

// 自动建表迁移
db.AutoMigrate(&User{},&Post{},&Comment{},&Like{})

//初始化gin
r:=gin.Default()

//注册路由
//公开接口
r.POST("/register",Register)
r.POST("/login",Login)

//登录接口
auth:=r.Group("/api")
auth.Use(AuthMiddleware())
{
	auth.POST("/posts",CreatePost)//f创建帖子
	auth.GET("/posts",ListPosts) //帖子列表
	auth.GET("/posts/:id",GetPost) //帖子详情
	auth.POST("/posts/:id/like",LikePost) //点赞帖子
	auth.DELETE("/posts/:id/like",UnlikePost) //取消点赞
	auth.POST("/posts/:id/comments",CreateComment)//评论
	auth.GET("/posts/:id/comments",ListComments) //评论区列表
}

	log.Println("服务器启动:8080")
	r.Run(":8080")

}
// 数据模型
type User struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Username  string  `gorm:"unique;not null" json:"用户名"`
	Password  string  `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"创建时间"`
}
type Post struct {
	ID  uint  `gorm:"primaryKey" json:"id"`
	Title string `gorm:"not null" json:"标题"`
	Content string `gorm:"not null" json:"内容"`
	UserID uint `gorm:"not null" json:"用户id"`
	User User `gorm:"foreignKey:UserID" json:"用户,信息"`
	CreatedAt time.Time `json:"创建时间"`
	LikeCount int `gorm:"default:0" json:"点赞数"`
}

type Comment struct {
	ID uint `gorm:"primarykey" json:"id"`
	Content string `gorm:"not null" json:"内容"`
	UserID uint `gorm:"not null" json:"user_id"`
	User User `gorm:"foreignkey:UserID" json:"用户信息"`
	PostID uint `gorm:"not null" json:"post_id"`
	CreatedAt time.Time `json:"创建时间"`
}
type Like struct{
	UserID uint `gorm:"primarykey" json:"用户id"`
	PostID uint `gorm:"primarykey" json:"帖子id"`
	CreatedAt time.Time `json:"点赞时间"`
} 
// jwt
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
// jwt验证
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少token"})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token格式错误"})
			c.Abort()
			return
		}
		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效token"})
					c.Abort()
					return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
// 接口处理函数
// 注册
func Register(c *gin.Context){
	var req struct {
		Username string ``
		Password string ``
	}
	if err:=c.ShouldBindJSON(&req);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return 
	}
	// 检查用户名是否重复
	var exist User
	if err:=db.Where("username=?",req.Username).First(&exist).Error;err==nil{
		c.JSON(http.StatusConflict,gin.H{"error":"用户名已存在"})
		return 
	}
	//加密密码
	 hashed,_:=bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)
	 user:=User{Username: req.Username,Password: string(hashed)}
	 db.Create(&user)
	 c.JSON(http.StatusOK,gin.H{"message":"注册成功","user_id":user.ID})
}
// 登录
	func Login (c*gin.Context){
		var req struct{
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err:=c.ShouldBindJSON(&req);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		var user User
		if err:=db.Where("username=?",req.Username).First(&user).Error;err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"用户密码错误"})
			return
		}
		if err:=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(req.Password));err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"用户密码错误"})
			return
		}
		token,_:=GenerateToken(user.ID,user.Username)
		c.JSON(http.StatusOK,gin.H{"token":token})
	}
	// 发帖
	func CreatePost (c *gin.Context){
		userID:=c.GetUint("userID")
		var req struct{
			Title string `json:"title" binding:"required"`
			Content string `json:"content" binding:"required"`
		}
		if err:=c.ShouldBindJSON(&req);err!=nil{
				c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
				return 
			}
			post :=Post{Title:req.Title,Content:req.Content,UserID:userID}
			db.Create(&post)
			c.JSON(http.StatusOK,post)
		
	}
	// 帖子列表
	func ListPosts(c *gin.Context){
		page,_:=strconv.Atoi(c.DefaultQuery("page","1"))
		pageSize,_:=strconv.Atoi(c.DefaultQuery("page_size", "10"))
		offset:=(page-1)*pageSize
		
		var posts []Post
		db.Preload("User").Order("created_at desc").Limit(pageSize).Offset(offset).Find(&posts)

		c.JSON(http.StatusOK,gin.H{"data":posts})
	}
	// 帖子详细
	func GetPost(c *gin.Context){
		id :=c.Param("id")
		var post Post
		if err:=db.Preload("User").First(&post,id).Error;err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"帖子不存在"})
			return
		}
		rdb.Incr(ctx,"post_view:" + id)
		c.JSON(http.StatusOK,post)
	}
	// 点赞
	func LikePost(c *gin.Context){
		userID:=c.GetUint("userID")
		postID,_:=strconv.Atoi(c.Param("id"))
		// 检查帖子是否存在
		var post Post
		if err :=db.First(&post,postID).Error;err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"帖子不存在"})
			return
		}
		//是否已点赞
		var like Like 
		err:=db.Where("user_id=? and post_id=?",userID,postID).First(&like).Error
		if err==nil{
			c.JSON(http.StatusContinue,gin.H{"error":"已经点过赞了"})
			return
		}
		// 事务;创建点赞记录+增加帖子的like_count
		err = db.Transaction(func(tx *gorm.DB)error{
			if err:=tx.Create(&Like{UserID:userID,PostID:uint(postID)}).Error;err!=nil{
				return err
			}
			if err:=tx.Model(post).Update("like_count",gorm.Expr("like_count+?,1")).Error;err!=nil{
				return err
			}
			return nil
		})
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"点赞失败"})
				return 
		}
		rdb.SAdd(ctx,"post_like:"+c.Param("id"),userID)
		c.JSON(http.StatusOK,gin.H{"massage":"点赞成功"})
	}
	// 取消点赞
	func UnlikePost(c *gin.Context){
		userID:=c.GetUint("userID")
		postID,_:=strconv.Atoi(c.Param("id"))
		// 事务
		err:=db.Transaction(func (tx*gorm.DB)error{
			// 删除点赞记录
			if err:=tx.Where("user_id=? and post_id=?",userID,postID).Delete(&Like{}).Error;err!=nil{
				return err
			}
			// 减少like_count(保证不小于0)
			if err:=tx.Model(&Post{}).Where("id=?",postID).Update("like_count",gorm.Expr("greatest(like_count-1,0)")).Error;err!=nil{
				return err 
			}
			return nil 
		})
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"errror":"取消失败"})
			return 
		}
		rdb.SRem(ctx,"post_likes:"+c.Param("id"),userID)
		c.JSON(http.StatusOK,gin.H{"massage":"取消点赞成功"})
	}
	// 评论列表o
	func ListComments(c *gin.Context){
		postID :=c.Param("id")
		var comments []Comment	
		db.Preload("User").Where("post_id=?",postID).Order("created_at desc").Find(&comments)
		c.JSON(http.StatusOK,comments)
	}
	// 发表评论
	func CreateComment(c *gin.Context){
		userID:=c.GetUint("userID")
		postID,err:=strconv.Atoi(c.Param("id"))
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"invalid post id"})
			return
		}

		var req struct{
			Content string `json:"content" binding:"required"`
		}
		if err :=c.ShouldBindJSON(&req);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return 
		}
		// 检查帖子是否存在、
		var post Post	
		if err:=db.First(&post,postID).Error;err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"帖子不存在"})
			return 
		}
		comment:=Comment{
			Content:req.Content,
			UserID:userID,
			PostID:uint(postID),
		}
		db.Create(&comment)
		c.JSON(http.StatusOK,comment)
	}