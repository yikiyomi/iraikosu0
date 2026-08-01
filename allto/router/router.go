package router

import (
    "allto/handler"
    "allto/middleware"
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	
	//注册路由
	r:=gin.New()

	//检测是否存货
	 r.GET("/healthz", func(c *gin.Context) {
      c.JSON(200, gin.H{"status": "ok"})
  	})
	
	//加载html
	r.Static("/static", "./frontend")
	//中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())
    r.Use(middleware.Logger())
    r.Use(middleware.RateLimit(60))

	//公开接口
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	//登录接口
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/posts", handler.CreatePost)                 //创建帖子
		auth.GET("/posts", handler.ListPosts)                   //帖子列表
		auth.GET("/posts/:id", handler.GetPost)                 //帖子详情
		auth.POST("/posts/:id/like", handler.LikePost)          //点赞帖子
		auth.DELETE("/posts/:id/like", handler.UnlikePost)      //取消点赞
		auth.POST("/posts/:id/comments", handler.CreateComment) //评论
		auth.GET("/posts/:id/comments", handler.ListComments)   //评论区列表
		auth.GET("/posts/:id/likes", handler.ListPostLikes)     //获取帖子的点赞列表
		auth.GET("/posts/likes", handler.ListMyLikes)           //当前用户点赞的帖子列表
		auth.GET("/posts/:id/likers", handler.ListPostLikers)   //帖子点赞用户列表
		auth.POST("/follow/:id", handler.FollowUser)            //关注用户
		auth.DELETE("/follow/:id", handler.UnfollowUser)        //取消关注
		auth.GET("/following", handler.Listfollowing)           //查看关注列表
		auth.GET("/followers", handler.GetFollowers)            //查看粉丝列表
	}
	return r
}
