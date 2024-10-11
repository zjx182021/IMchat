package service

import (
	"TM_chat/models"
	"TM_chat/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/rand"
)

// /user/getUserList GET
func GetUserList(ctx *gin.Context) {
	// data := make([]*models.UserBasic, 10)
	data := models.GetUserList()
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "用户已注册",
		"data":    data,
	})
}

// /user/createUser 处理创建新用户的请求 PUT
func Createuser(ctx *gin.Context) {
	// 初始化一个新的 UserBasic 对象
	user := &models.UserBasic{}

	// 从查询参数中获取用户的名字
	user.Name = ctx.PostForm("name")
	// 从查询参数中获取用户的密码
	password := ctx.PostForm("password")
	// 从查询参数中获取用户的重复密码
	repassword := ctx.PostForm("Identity")
	// 如果两次输入的密码不一致，则返回错误信息
	if password != repassword {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "两次密码不一致",
			"data":    user,
		})
		return
	}
	user.Salt = fmt.Sprintf("%06d", rand.Int31())
	fmt.Println(password)
	// 将密码设置给用户对象
	user.Password = utils.MakePassword(password, user.Salt)
	fmt.Println(user.Password)
	// 调用 models 包中的 CreateUser 函数来创建用户
	models.CreateUser(user)
	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    user,
	})
}

// /user/deleteUser 删除 GET
func DeleteUser(ctx *gin.Context) {
	user := &models.UserBasic{}
	// 从查询参数中获取用户的名字
	id, _ := strconv.Atoi(ctx.Query("Id"))

	user.ID = uint(id)
	models.DeleteUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"id":      user.ID,
	})
}

// /user/updateUser 更新 POST
func UpdateUser(ctx *gin.Context) {
	user := &models.UserBasic{}
	id, _ := strconv.Atoi(ctx.PostForm("Id"))
	user.ID = uint(id)
	user.Name = ctx.DefaultPostForm("Name", user.Name)
	user.Password = ctx.DefaultPostForm("Password", user.Password)
	user.Phone = ctx.DefaultPostForm("Phone", user.Phone)
	email := ctx.PostForm("Email")
	if email != "" {
		user.Email = &email
	}
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return // 验证失败，返回并结束函数执行
	}
	models.UpdateUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"message":   "更新成功",
		"Name":      user.Name,
		"Id":        user.ID,
		"Phone":     user.Phone,
		"Email":     user.Email,
		"PassworD":  user.Password,
		"LoginTime": user.LoginTime,
	})
}

// r.POST("/user/Login", service.LoginCheck)
func LoginCheck(ctx *gin.Context) {

	name := ctx.PostForm("name")
	pwd := ctx.PostForm("password")
	user := models.UserBasic{Name: name}
	_ = models.FindUserByname(&user)
	fmt.Println(user)
	identity := ctx.DefaultPostForm("identity", "0")
	if user.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "用户不存在",
		})
		return
	}

	if user.Identity == "" {
		if utils.ValidPassword(pwd, user.Salt, user.Password) {
			_ = models.Updateidentity(&user)
			fmt.Println(user.Identity)
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "登陆成功",
				"data":    user,
				"token":   user.Identity,
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "密码不正确",
			})
		}
	} else {
		if identity == user.Identity {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "登陆成功",
				"data":    user,
				"token":   user.Identity,
			})
		} else {
			if utils.ValidPassword(pwd, user.Salt, user.Password) {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    0,
					"message": "登陆成功",
					"data":    user,
				})
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    -1,
					"message": "身份验证不通过",
				})
			}
		}
	}
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(ctx *gin.Context) {
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println("close:", err)
		}
	}(ws)
	MsgHandler(ws, ctx)
}

func MsgHandler(ws *websocket.Conn, ctx *gin.Context) {
	msg, err := utils.Subscribe(ctx, utils.PublishKey)
	if err != nil {
		fmt.Println("subscribe:", err)
		return
	}
	tm := time.Now().Format("2001-01-01 01:01:01")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println("write:", err)
		return
	}
}

// r.GET("/user/sendUserMsg", service.SendUserMsg)
func SendUserMsg(ctx *gin.Context) {
	models.Chat(ctx.Writer, ctx.Request)
}

func SearchFriends(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("userId"))
	fmt.Println(id)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "userId error")
		return
	}
	user := models.SearchFriend(uint(id))
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"code":    0,
	// 	"message": "查找成功",
	// 	"data":    user,
	// })
	utils.RespOKList(ctx.Writer, user, len(user))
}
