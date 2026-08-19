package handler

import (
	"allto/database"
	"allto/model"
	"allto/response"
	"allto/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 帖子的点赞用户列表
func ListPostLikes(c *gin.Context) {
	postID := c.Param("id")
	_, err := strconv.Atoi(postID)
	if err != nil {
		response.BadRequest(c, "无效帖子id")
		return
	}
	var likes []model.Like
	if err := database.GetDB().Where("post_id = ?", postID).Order("created_at desc").Find(&likes).Error; err != nil {

		response.InternalError(c, "查询点赞的用户列表失败")
		return
	}
	response.Success(c, likes)
}

// 获取帖子的点赞用户列表
func ListPostLikers(c *gin.Context) {
	postID := c.Param("id")
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		response.BadRequest(c, "无效帖子id")
		return
	}

	// 使用JOIN查询关联表
	var likers []model.User
	if err := database.GetDB().Joins("JOIN likes ON users.id = likes.user_id").
		Where("likes.post_id = ?", postIDInt).
		Select("users.id, users.username").
		Find(&likers).Error; err != nil {
		response.InternalError(c, "查询失败")
		return
	}
	response.Success(c, likers)
}

// 点赞
func LikePost(c *gin.Context) {
	userID := c.GetUint("userID")
	postID, _ := strconv.Atoi(c.Param("id"))
	// 检查帖子是否存在
	var post model.Post
	if err := database.GetDB().First(&post, postID).Error; err != nil {
		response.NotFound(c, "帖子不存在")
		return
	}
	//是否已点赞
	var like model.Like
	err := database.GetDB().Where("user_id=? and post_id=?", userID, postID).First(&like).Error
	if err == nil {
		response.BadRequest(c, "已经点赞了")
		return
	}
	// 事务;创建点赞记录+增加帖子的like_count
	err = database.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model.Like{UserID: userID, PostID: uint(postID)}).Error; err != nil {
			return err
		}
		if err := tx.Model(&post).Update("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		response.InternalError(c, "点赞失败")
		return
	}
	err=database.SAdd("post_like:"+c.Param("id"),userID)//reids歇逼放行
	if err!=nil{
		util.Logger.Error("点赞缓存写入失败(降级)",zap.Error(err))
	}
	go Notify(post.UserID, userID, "like", uint(postID))
	response.Success(c, "点赞成功")
}

// 取消点赞
func UnlikePost(c *gin.Context) {
	userID := c.GetUint("userID")
	postID, _ := strconv.Atoi(c.Param("id"))
	// 事务
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 删除点赞记录
		if err := tx.Where("user_id=? and post_id=?", userID, postID).Delete(&model.Like{}).Error; err != nil {
			return err
		}
		// 减少like_count(保证不小于0)
		if err := tx.Model(&model.Post{}).Where("id=?", postID).Update("like_count", gorm.Expr("greatest(like_count-1,0)")).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		response.InternalError(c, "取消失败")
		return
	}
	err=database.SRem("post_like:"+c.Param("id"), userID)//reids歇逼放行
	if err!=nil{
		util.Logger.Error("取消点赞缓存删除失败(降级)",zap.Error(err))
	}
	response.Success(c, "取消成功")
}

// 评论列表o
func ListComments(c *gin.Context) {
	postID := c.Param("id")
	var comments []model.Comment
	database.GetDB().Preload("User").Where("post_id=?", postID).Order("created_at desc").Find(&comments)
	response.Success(c, comments)
}

// 发表评论
func CreateComment(c *gin.Context) {
	userID := c.GetUint("userID")
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	// 检查帖子是否存在、
	var post model.Post
	if err := database.GetDB().First(&post, postID).Error; err != nil {
		response.NotFound(c, "帖子不存在")
		return
	}
	comment := model.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  uint(postID),
	}
	database.GetDB().Create(&comment)
	go Notify(post.UserID, userID, "comment", uint(postID))
	response.Success(c, comment)
}
// 当前用户点赞列表
func ListMyLikes(c *gin.Context) {
	userID := c.GetUint("userID")
	var likes []model.Like
	if err := database.GetDB().Preload("Post").Where("user_id=?", userID).Order("created_at desc").Find(&likes).Error; err != nil {
		response.InternalError(c, "查询失败")
		return
	}
	response.Success(c, likes)
}