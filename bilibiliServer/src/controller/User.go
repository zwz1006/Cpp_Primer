package controller

import (
	"bilibiliServer/src/model"
	"bilibiliServer/src/tools"
	"log"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(c *gin.Context) { // 用户注册控制
	db := tools.GetDb()

	var user model.User
	c.BindJSON(&user)
	if b, _ := regexp.MatchString("^[^ ]{2,7}$", user.Name); !b {
		c.JSON(400, gin.H{"msg": "昵称不符合规范!"})
		return
	}
	if b, _ := regexp.MatchString("^[0-9]{6,11}$", user.Account); !b {
		c.JSON(400, gin.H{"msg": "用户名不符合规范!"})
		return
	}
	if b, _ := regexp.MatchString("^[a-zA-Z0-9_-]{6,11}$", user.Password); !b {
		c.JSON(400, gin.H{"msg": "密码不符合规范!"})
		return
	}
	if isUserExits(db, user.Account) {
		c.JSON(400, gin.H{"msg": "该用户帐号已存在!"})
	} else {
		db.Create(&user)
		token := tools.GetToken(strconv.Itoa(int(user.ID))) //  给注册成功的用户返回 token 和 id
		c.JSON(200, gin.H{"msg": "注册成功~", "id": user.ID, "token": token})
	}

}
func Login(c *gin.Context) { // 用户登录
	db := tools.GetDb()
	var remoteUser model.User
	c.BindJSON(&remoteUser)
	if b, _ := regexp.MatchString("^[0-9]{6,11}$", remoteUser.Account); !b {
		c.JSON(400, gin.H{"msg": "用户名不规范"})
		return
	}
	if b, _ := regexp.MatchString("^[a-zA-Z0-9_-]{6,11}$", remoteUser.Password); !b {
		c.JSON(400, gin.H{"msg": "密码不符合规范!"})
		return
	}

	if isUserExits(db, remoteUser.Account) {
		var localUser model.User
		db.Where("Account = ?", remoteUser.Account).First(&localUser)
		if remoteUser.Password != localUser.Password {
			c.JSON(401, gin.H{"msg": "密码错误!"})
		} else {
			token := tools.GetToken(strconv.Itoa(int(localUser.ID))) // 给登录成功的用户返回 id 和 token
			c.JSON(200, gin.H{"msg": "登陆成功~", "id": localUser.ID, "token": token})
		}
	} else {
		c.JSON(400, gin.H{"msg": "帐号不存在!"})
	}
}

func UserInfo(c *gin.Context) { // 用户信息查询
	id := c.Param("id")

	log.Println("用户 " + id + " 进行了一次数据查询")
	db := tools.GetDb()
	var user model.User
	db.Where("id = ?", id).First(&user)
	c.JSON(200, gin.H{"name": user.Name, "following": user.Following, "followers": user.Followers, "likes": user.Likes, "introduction": user.Introduction, "account": user.Account, "avatar": tools.AvatarPath(user.Avatar), "sex": user.Sex})
}

func UserModify(c *gin.Context) { // 用户资料修改
	id := c.Param("id")
	log.Println("用户: ", id, " 正在修改资料")
	var user model.User
	c.BindJSON(&user)

	db := tools.GetDb()
	if user.Name != "" {
		if b, _ := regexp.MatchString("^[^ ]{2,7}$", user.Name); !b {
			c.JSON(400, gin.H{"msg": "昵称不符合规范!"})
			return
		} // 有用户名且符合规范, 修改用户名
		db.Model(&user).Where("id = ?", id).Update("name", user.Name)
		log.Println("用户", id, "修改用户名成功~")
		c.JSON(200, gin.H{"msg": "用户名修改成功~"})
	}
	if user.Introduction != "" {
		// 个性签名不为空, 修改个性签名
		db.Model(user).Where("id = ?", id).Update("introduction", user.Introduction)
		log.Println("用户", id, "修改签名成功~")

		c.JSON(200, gin.H{"msg": "个性签名修改成功~"})
	}
	if user.Sex != "" {
		if user.Sex != "男" && user.Sex != "女" {
			c.JSON(400, gin.H{"msg": "性别不符合规范!"})
			return
		}
		// 性别不为空且符合规范, 修改性别
		db.Model(user).Where("id = ?", id).Update("sex", user.Sex)
		log.Println("用户", id, "修改性别成功~")

		c.JSON(200, gin.H{"msg": "性别修改成功~"})
	}
}

func isUserExits(db *gorm.DB, account string) bool {
	var user model.User
	db.Where("account = ?", account).First(&user)
	if user.ID != 0 {
		return true
	} else {
		return false
	}
}