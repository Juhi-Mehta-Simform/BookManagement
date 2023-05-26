package router

import (
	"Project/controller"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	Server = gin.Default()
)

func init() {
	store := cookie.NewStore([]byte("secret"))
	Server.Use(sessions.Sessions("mysession", store))
	Server.LoadHTMLGlob("templates/*")
	Server.GET("/", controller.LoadLogin)
	Server.POST("/login", controller.Login)
	Server.GET("/logout", controller.Logout)
	Server.GET("/load", controller.LoadSignup)
	Server.POST("/signup", controller.Signup)
	Server.GET("/addBook", controller.LoadBook)
	Server.POST("/addBook", controller.CreateBook)
	Server.GET("/viewBook", controller.ViewBook)
	Server.GET("/loadBook/:id", controller.LoadUpdate)
	Server.POST("/updateBook", controller.UpdateBook)
	Server.GET("/loadDelete/:id", controller.DeleteBook)
	Server.GET("/loadBorrow/:id", controller.LoadBorrow)
	Server.POST("/borrowBook", controller.BorrowBook)
	Server.GET("/viewBorrow", controller.ViewBorrow)
	Server.GET("/userBorrow", controller.UserBorrow)
	Server.GET("/loadReturn/:id", controller.LoadReturn)
}
