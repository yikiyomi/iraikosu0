package handler

import (
	"allto/database"
	"allto/model"
	"allto/response"
	"strconv"
	"log"
	"github.com/gin-gonic/gin"
)

// 发帖
func CreatePost(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	post := model.Post{Title: req.Title, Content: req.Content, UserID: userID}
	database.GetDB().Create(&post)
	response.Success(c, post)

}

// 帖子列表
func ListPosts(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		response.BadRequest(c, "page 页码解析错误")
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil {
		response.BadRequest(c, "pagesize 页码解析错误")
		return
	}
	if pageSize > 50 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	var posts []model.Post
	database.GetDB().Preload("User").Order("created_at desc").Limit(pageSize).Offset(offset).Find(&posts)
	response.Success(c, posts)
}

// 帖子详细
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post model.Post
	if err := database.GetDB().Preload("User").First(&post, id).Error; err != nil {
		response.NotFound(c, "帖子不存在")
		return
	}
	err := database.SafeIncr("post_view:" + id)//reids歇逼放行
	if err != nil {
		log.Printf("阅读数统计失败(降级): %v", err)
	}
	response.Success(c, post)
}

