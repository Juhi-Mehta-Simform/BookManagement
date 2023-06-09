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
	Server.LoadHTMLGlob("templates/html/*")
	Server.Static("/css", "templates/css")
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
	Server.GET("/loadBorrow/:id/:isbn", controller.LoadBorrow)
	Server.POST("/borrowBook", controller.BorrowBook)
	Server.GET("/viewBorrow", controller.ViewBorrow)
	Server.GET("/viewRequest", controller.ViewRequest)
	Server.GET("/userBorrow", controller.UserBorrow)
	Server.GET("/loadReturn/:book_isbn", controller.LoadReturn)
	Server.GET("/returnRequest/:member_id/:book_isbn", controller.ReturnRequest)
	Server.POST("/returnBook", controller.ReturnBook)
	Server.GET("/makeLibrarian/:user_id", controller.MakeLibrarian)
	Server.GET("/removeLibrarian/:user_id", controller.RemoveLibrarian)
	Server.GET("/makeAdmin/:user_id", controller.MakeAdmin)
	Server.GET("/fetch", controller.GetUser)
	Server.GET("/searchFilterBook", controller.SearchFilterBook)
	Server.GET("/viewProfile", controller.UserProfile)
	Server.GET("/loadProfile/:user_id", controller.LoadProfile)
	Server.POST("/updateProfile", controller.UpdateProfile)
	Server.GET("/forgetPassword", controller.LoadForget)
	Server.POST("/forgetPassword", controller.ForgetPassword)
	Server.GET("/resetPassword", controller.LoadReset)
	Server.POST("/resetPassword", controller.ResetPassword)
	Server.GET("/loadDonate", controller.LoadDonate)
	Server.POST("/donateBook", controller.DonateBook)
	Server.GET("/userDonate", controller.UserDonate)
	Server.GET("/viewDonate", controller.ViewDonate)
	Server.GET("/searchFilterUser", controller.SearchFilterUser)
	Server.GET("/borrowHistory/:user_id", controller.BorrowHistory)
	Server.GET("/donateHistory/:user_id", controller.DonateHistory)
	Server.GET("/reminder/:member_id", controller.Reminder)
}
