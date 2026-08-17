package handler
import (
	"allto/database"
	"allto/model"
	"allto/response"
	"allto/util"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
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