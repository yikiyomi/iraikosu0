package handler

import (
	"allto/database"
	"allto/model"
	"allto/response"
	"allto/util"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 接口处理函数
// 注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	// 检查用户名是否重复
	var exist model.User
	if err := database.GetDB().Where("username=?", req.Username).First(&exist).Error; err == nil {
		response.Conflict(c, "用户名已存在")
		return
	}
	verifyToken,err:=util.GenerateRefreshToken(16)
	if err!=nil{
		response.InternalError(c,"生成验证token失败")
		return
	}
	//检查邮箱是否占用
	if err := database.GetDB().Where("email = ?", req.Email).First(&exist).Error; err == nil {
        response.Conflict(c, "邮箱已被注册")
        return
    }
	//加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败")
		response.InternalError(c, "注册失败,请稍后再试")
		return
	}
	user := model.User{Username: req.Username, Password: string(hashed),Email: req.Email,VerifyToken: verifyToken,}
	database.GetDB().Create(&user)
	verifyLink := "http://localhost:8080/verify-email?token=" + verifyToken
  	go util.SendVerifyEmail(req.Email, verifyLink)
	response.Success(c, gin.H{"message": "注册成功", "user_id": user.ID})
}

// 登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var user model.User
	if err := database.GetDB().Where("username=?", req.Username).First(&user).Error; err != nil {
		response.Unauthorized(c, "账号不存在")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "密码错误")
		return
	}
	if !user.EmailVerified && user.Email != "" {
    response.Unauthorized(c, "请先验证邮箱")
    return
    }
	//霜双token并存入数据库
	token, err := util.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.InternalError(c, "token生成失败")
		return
	}
	refreshToken, err := util.GenerateRefreshToken(32)
	if err != nil {
		response.InternalError(c, "refresh token生成失败")
		return
	}
	result := database.GetDB().Model(&user).Updates(map[string]interface{}{
		"token":         token,        // 当前 access_token
		"refresh_token": refreshToken, // 当前 refresh_token
	})
	if result.Error != nil {
		log.Printf("登陆时的双token存入数据库失败")
		response.InternalError(c, "token存入数据库失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"refresh_token": refreshToken,
		"user":          gin.H{
			"id": user.ID,
			"username": user.Username,
			"email": user.Email,                 
        	"email_verified": user.EmailVerified,}})
}

// 登出
func Logout(c *gin.Context) {
	userID := c.GetUint("userID")
	database.GetDB().Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"token":         "",
		"refresh_token": "",
	})
	response.Success(c, "已登出")
}

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
	database.GetRedis().Incr(database.GetCtx(), "post_view:"+id)
	response.Success(c, post)
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
	database.GetRedis().SAdd(database.GetCtx(), "post_like:"+c.Param("id"), userID)
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
	database.GetRedis().SRem(database.GetCtx(), "post_like:"+c.Param("id"), userID)
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
	response.Success(c, comment)
}

// 关注用户
func FollowUser(c *gin.Context) {
	followID := c.GetUint("userID")
	followingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效用户id")
		return
	}
	//不能关注自己
	if followID == uint(followingID) {
		response.BadRequest(c, "不能关注自己")
		return
	}
	// 检查关注d用户是否存在
	var user model.User
	if err := database.GetDB().First(&user, followingID).Error; err != nil {
		response.NotFound(c, "关注的用户不存在")
		return
	}
	//是否已关注
	var follow model.Follow
	err = database.GetDB().Where("follow_id=? and following_id=?", followID, followingID).First(&follow).Error
	if err == nil {
		response.BadRequest(c, "已关注")
		return
	}
	// 事务;创建关注记录
	nfollow := model.Follow{FollowID: followID, FollowingID: uint(followingID)}
	database.GetDB().Create(&nfollow)
	response.Success(c, gin.H{"message": "关注成功"})
}

// 取消关注
func UnfollowUser(c *gin.Context) {
	followID := c.GetUint("userID")
	followingID, _ := strconv.Atoi(c.Param("id"))

	result := database.GetDB().Where("follow_id= ? and following_id=?", followID, followingID).Delete(&model.Follow{})
	if result.RowsAffected == 0 {
		response.BadRequest(c, "未关注该用户")
		return
	}
	response.Success(c, "取消成功")
}

// 查看关注列表
func Listfollowing(c *gin.Context) {
	userID := c.GetUint("userID")

	var follows []model.Follow
	if err := database.GetDB().Where("follow_id = ? ", userID).Find(&follows).Error; err != nil {
		response.InternalError(c, "查询关注列表失败")
		return
	}

	//提取所有被关注者id
	var ids []uint
	for _, f := range follows {
		ids = append(ids, f.FollowingID)
	}

	var users []model.User
	if len(ids) > 0 {
		if err := database.GetDB().Where("id in ?", ids).Find(&users).Error; err != nil {
			response.InternalError(c, "查询失败")
			return
		}
		response.Success(c, users)
	} else {
		response.Success(c, []model.User{})
	}
}

// 查看粉丝列表
func GetFollowers(c *gin.Context) {
	userID := c.GetUint("userID")

	var follows []model.Follow
	if err := database.GetDB().Where("following_id = ? ", userID).Find(&follows).Error; err != nil {
		response.InternalError(c, "查询粉丝列表失败")
		return
	}

	//提取所有关注者id
	var ids []uint
	for _, f := range follows {
		ids = append(ids, f.FollowID)
	}

	var users []model.User
	if len(ids) > 0 {
		if err := database.GetDB().Where("id in ?", ids).Find(&users).Error; err != nil {
			response.InternalError(c, "查询粉丝列表失败")
			return
		}
		response.Success(c, users)
	} else {
		response.Success(c, []model.User{})
	}
}

func UserProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效用户id")
		return
	}
	var user model.User
	if database.GetDB().First(&user, userID).Error != nil {
		response.NotFound(c, "找不到该用户")
		return
	}
	var postCount int64
	var followerCount int64
	var followingCount int64

	database.GetDB().Model(&model.Post{}).Where("user_id = ?", userID).Count(&postCount)            //贴子数
	database.GetDB().Model(&model.Follow{}).Where("following_id = ?", userID).Count(&followerCount) // 粉丝
	database.GetDB().Model(&model.Follow{}).Where("follow_id = ?", userID).Count(&followingCount)   // 关注

	//打包返回
	response.Success(c, gin.H{
		"id":              user.ID,
		"username":        user.Username,
		"created_at":      user.CreatedAt,
		"post_count":      postCount,
		"bio":             user.Bio,
		"follower_count":  followerCount,
		"following_count": followingCount,
		"avatar_url":      user.AvatarUrl,
	})
}

func UserPosts(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效用户id")
		return
	}
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
	database.GetDB().Preload("User").
		Where("user_id = ?", userID). //只查该作者
		Order("created_at desc").
		Limit(pageSize).Offset(offset).
		Find(&posts)

	response.Success(c, posts)
}
func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		Bio string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	database.GetDB().Model(&model.User{}).Where("id=?", userID).Update("bio", req.Bio)
	response.Success(c, "修改成功")
}

// 头像
func UploadAvatar(c *gin.Context) {
	userID := c.GetUint("userID")
	//获取文件
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	allowedExts := map[string]bool{".jpg": true, ".png": true}
	if !allowedExts[ext] {
		response.BadRequest(c, "头像仅支持jpg,png")
		return
	}
	if header.Size > 5*1024*1024 {
		response.BadRequest(c, "头像大小不超过5MB")
		return
	}

	filename := fmt.Sprintf("%s%s", uuid.NewString(), ext)
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalError(c, "服务器文件系统错误")
		return
	}

	destPath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(destPath)
	if err != nil {
		response.InternalError(c, "保存头像失败")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.InternalError(c, "写入头像失败")
		return
	}

	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)
	if err := database.GetDB().Model(&model.User{}).Where("id=?", userID).Update("avatar_url",
		avatarURL).Error; err != nil {
		os.Remove(destPath)
		response.InternalError(c, "更新头像记录失败")
		return
	}
	response.Success(c, gin.H{"avatar_url": avatarURL})
}

// refresh token
func Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var user model.User
	if err := database.GetDB().Where("refresh_token=?", req.RefreshToken).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Unauthorized(c, "refresh token无效或过期")
		} else {
			response.InternalError(c, "服务端错误")
		}
		return
	}
	newAccessToken, err := util.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.InternalError(c, "生成 access token 失败")
		return
	}

	// 3. 生成新的 refresh_token（滚动更新）
	newRefreshToken, err := util.GenerateRefreshToken(16)
	if err != nil {
		response.InternalError(c, "生成 refresh token 失败")
		return
	}
	updates := map[string]interface{}{
		"token":         newAccessToken,
		"refresh_token": newRefreshToken,
	}
	if err := database.GetDB().Model(&model.User{}).Where("id=?", user.ID).Updates(updates).Error; err != nil {
		response.InternalError(c, "更新token失败")
		return
	}
	response.Success(c, gin.H{
		"token":         newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// 改密码
func ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		response.BadRequest(c, "用户不存在")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		response.Unauthorized(c, "旧密码错误")
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}
	err = database.GetDB().Model(&user).Updates(map[string]interface{}{
		"password":      string(hashed),
		"token":         "",
		"refresh_token": "",
	}).Error
	if err != nil {
		response.InternalError(c, "密码更新失败")
		return
	}
	response.Success(c, "修改成功，请重新登陆")
}

// 改名
func Rename(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		NewUsername string `json:"new_username" binding:"required"`
	}
	var existing model.User
	err := database.GetDB().Where("username = ? and id != ?", req.NewUsername, userID).First(&existing).Error
	if err == nil {
		response.Conflict(c, "用户名已被占用")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		response.InternalError(c, "查询用户失败")
		return
	}
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		response.Unauthorized(c, "用户不存在")
		return
	}
	if err := database.GetDB().Model(&user).Update("username", req.NewUsername).Error; err != nil {
		response.InternalError(c, "更新用户名失败")
		return
	}
	newToken, err := util.GenerateToken(user.ID, req.NewUsername)
	if err != nil {
		response.InternalError(c, "otken生成失败")
		return
	}
	if err := database.GetDB().Model(&user).Update("token", newToken).Error; err != nil {
		response.InternalError(c, "更新 token 失败")
		return
	}
	response.Success(c, gin.H{
		"token":    newToken,
		"username": req.NewUsername,
	})
}
func VerifyEmail(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")

	token := c.Query("token")
	if token == "" {
		c.String(http.StatusBadRequest, renderVerifyPage("验证链接无效", "链接缺少 token 参数"))
		return
	}

	var user model.User
	if err := database.GetDB().Where("verify_token = ?", token).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, renderVerifyPage("验证链接无效或已过期", "请重新绑定邮箱获取新的验证链接"))
		} else {
			c.String(http.StatusInternalServerError, renderVerifyPage("服务器错误", "请稍后再试"))
		}
		return
	}

	if err := database.GetDB().Model(&user).Updates(map[string]interface{}{
		"email_verified": true,
		"verify_token":   "",
	}).Error; err != nil {
		c.String(http.StatusInternalServerError, renderVerifyPage("验证失败", "请稍后再试"))
		return
	}

	c.String(http.StatusOK, renderVerifyPage("✅ 验证成功", "你的邮箱已验证，现在可以登录了"))
}

// renderVerifyPage 生成邮箱验证结果页面
func renderVerifyPage(title, msg string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>%s</title>
<style>body{font-family:sans-serif;text-align:center;padding-top:80px;color:#333}h2{font-size:24px}p{color:#666}a{color:#1a73e8}</style>
</head><body><h2>%s</h2><p>%s</p><p><a href="/static/index.html">返回登录</a></p></body></html>`, title, title, msg)
}
func BindEmail(c *gin.Context){
	userID := c.GetUint("userID")
	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		response.BadRequest(c, "邮箱格式不正确")
		return
	}

	var user model.User
	err := database.GetDB().Where("email= ?", req.Email).First(&user).Error
	if err ==nil{
		response.NotFound(c,"邮箱被占用")
	} else if err != gorm.ErrRecordNotFound {
		response.InternalError(c, "系统繁忙，请稍后重试")
		return
	}

	// 生成验证 token
	verifyToken, err := util.GenerateRefreshToken(32)
	if err != nil {
		response.InternalError(c, "生成验证token失败")
		return
	}

	// 更新邮箱 + 验证 token，并重置验证状态
	if err := database.GetDB().Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"email":          req.Email,
		"verify_token":   verifyToken,
		"email_verified": false,
	}).Error; err != nil {
		response.InternalError(c, "绑定邮箱失败，请稍后重试")
		return
	}

	// 异步发送验证邮件，链接从请求 Host 推断（本地与部署都适用）
	go func() {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		verifyURL := scheme + "://" + c.Request.Host + "/verify-email?token=" + verifyToken
		if err := util.SendVerifyEmail(req.Email, verifyURL); err != nil {
			log.Printf("发送验证邮件失败: user_id=%d, email=%s, err=%v", userID, req.Email, err)
		}
	}()

	response.Success(c, "验证邮件已发送，请查收")
}