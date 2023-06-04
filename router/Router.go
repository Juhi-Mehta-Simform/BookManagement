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
	Server.GET("/home", controller.Home)
	Server.GET("/viewUser", controller.ViewUser)
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
	Server.GET("/loadReturn/:book_isbn", controller.LoadReturn)
	Server.GET("/returnRequest/:member_id/:book_isbn", controller.ReturnRequest)
	Server.POST("/returnBook", controller.ReturnBook)
	Server.GET("/makeLibrarian/:user_id", controller.MakeLibrarian)
	Server.GET("/makeAdmin/:user_id", controller.MakeAdmin)
	Server.POST("/search", controller.SearchUser)
	Server.POST("/searchBook", controller.SearchBook)
	Server.GET("/fetch", controller.GetUser)
	Server.GET("/filterBook", controller.FilterBook)
	Server.GET("/viewProfile", controller.UserProfile)
	Server.GET("/loadProfile/:user_id", controller.LoadProfile)
	Server.POST("/updateProfile", controller.UpdateProfile)
	Server.GET("/forgetPassword", controller.LoadForget)
	Server.POST("/forgetPassword", controller.ForgetPassword)
	Server.GET("/resetPassword/:email", controller.LoadReset)
	Server.POST("/resetPassword", controller.ResetPassword)
	Server.GET("/loadDonate", controller.LoadDonate)
	Server.POST("/donateBook", controller.DonateBook)
	Server.GET("/userDonate", controller.UserDonate)
	Server.GET("/viewDonate", controller.ViewDonate)
	Server.GET("/filterUser", controller.FilterUser)
}
