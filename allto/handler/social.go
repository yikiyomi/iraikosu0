package handler
import (
	"allto/database"
	"allto/model"
	"allto/response"
	"strconv"
	"github.com/gin-gonic/gin"
)
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