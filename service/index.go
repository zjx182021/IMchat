package service

import (
	"TM_chat/models"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// router /index
func GetIndex(ctx *gin.Context) {
	// 获取首页数据
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ind.Execute(ctx.Writer, "index")
}

func ToRegister(ctx *gin.Context) {
	// 获取首页数据
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ind.Execute(ctx.Writer, "register")
}

func ToChat(ctx *gin.Context) {
	// 获取首页数据
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	user := &models.UserBasic{}
	userId, err := strconv.Atoi(ctx.Query("userId"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "userId error")
		return
	}
	token := ctx.Query("token")
	user.ID = uint(userId)
	user.Identity = token
	ind.Execute(ctx.Writer, user)
}

func Chat(ctx *gin.Context) {
	models.Chat(ctx.Writer, ctx.Request)
}
