package handler

import (
	"allto/database"
	"allto/middleware"
	"allto/model"
	"allto/response"
	"allto/util"
	"fmt"
	"encoding/json"
	"github.com/gin-gonic/gin"
)
var hub=middleware.NewSSEHub()

// Notify 写通知 + 推 SSE。recipientID=接收者，senderID=触发者，typ=类型，targetID=目标
func Notify(recipientID, senderID uint, typ string, targetID uint) {
    if recipientID == senderID {
        return // 不给自己发通知
    }
    n := model.Notification{
        RecipientID: recipientID,
    	SenderID:    senderID,
        Type:        typ,
        TargetID:    targetID,
    }
    database.GetDB().Create(&n)

    // 推 SSE（失败也不影响，通知已存 DB）
    msg, _ := json.Marshal(gin.H{
      "type":     typ,
      "sender_id": senderID,
      "target_id": targetID,
    })
    hub.Push(recipientID, msg)
}
// NotificationStream SSE 长连接（token 放 query string，因为 EventSource 不能带 header）
func NotificationStream(c *gin.Context){
	token:=c.Query("token")
	claims,err:=util.ParseToken(token)
	if err!=nil{
		c.JSON(401,gin.H{"error":"无效token"})
		return
	}
	userID:=claims.UserID

	c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")

	ch:=hub.Register(userID)
	defer hub.Unregister(userID)
	c.Writer.Flush()

	for{
		select {
		case msg:=<-ch:
			fmt.Fprintf(c.Writer,"data: %s\n\n",msg)
			c.Writer.Flush()
		case<-c.Request.Context().Done():
			return
		}
	}
}
//通知列表（50条）
func ListNotifications(c *gin.Context){
	userID:=c.GetUint("userID")
	var list []model.Notification
	database.GetDB().Where("recipient_id=?",userID).
	Order("created_at desc").Limit(50).Find(&list)
	response.Success(c,list)
}
// UnreadCount 未读计数
  func UnreadCount(c *gin.Context) {
    userID := c.GetUint("userID")
    var count int64
    database.GetDB().Model(&model.Notification{}).
        Where("recipient_id = ? AND is_read = ?", userID, false).
        Count(&count)
    response.Success(c, gin.H{"count": count})
}
//标记已读（看id单条
func MarkRead(c *gin.Context){
	userID:=c.GetUint("userID")
	var req struct{
		ID uint `json:"id"`
	}
	c.ShouldBind(&req)

	q := database.GetDB().Model(&model.Notification{}).Where("recipient_id = ?", userID)
      if req.ID > 0 {
        q = q.Where("id = ?", req.ID)
    }
    q.Update("is_read", true)
    response.Success(c, "已读")
}