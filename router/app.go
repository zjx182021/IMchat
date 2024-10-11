package router

import (
	"TM_chat/models"
	"TM_chat/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func CreateBeforeCheckByname(ctx *gin.Context) {
	name := ctx.PostForm("name")
	user := &models.UserBasic{Name: name}
	_ = models.FindUserByname(user)
	if user.ID != 0 {
		ctx.JSON(200, gin.H{
			"message": "该名称已存在",
			"code":    -1,
			"data":    user,
		})
		ctx.Abort()
	}
	ctx.Next()
}

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Static("/asset", "asset/")
	r.LoadHTMLFiles("views/**/*")
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.POST("/searchFriends", service.SearchFriends)
	r.GET("/chat", service.Chat)
	r.GET("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", CreateBeforeCheckByname, service.Createuser)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/Login", service.LoginCheck)
	r.GET("/user/send", service.SendMsg)
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	r.POST("/attach/upload", service.Upload)
	// r.POST("/user/Register", service.ToRegister)
	r.POST("/contact/addfriend", service.AddFriend)
	return r
}
