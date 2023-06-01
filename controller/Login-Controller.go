package controller

import (
	"Project/connection"
	"Project/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func LoadSignup(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", nil)
}

func LoadLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func Signup(ctx *gin.Context) {
	var user models.User
	user.RoleName = ctx.PostForm("role_name")
	user.Name = ctx.PostForm("name")
	user.Email = ctx.PostForm("email")
	user.Password = ctx.PostForm("password")
	err := connection.GetConnection().Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"error": "User already exists",
			})
		}
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func Login(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user models.User
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	connection.GetConnection().Debug().Model(&models.User{}).Where("email=?", email).Find(&user)
	if user.UserID != 0 {
		db := connection.GetConnection().Debug().Model(&models.User{}).Where("email=? AND password=?", email, password).Find(&user)
		if db.RowsAffected == 0 {
			ctx.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error": "Incorrect Password",
			})
		} else {
			session.Set("userID", user.UserID)
			session.Save()
			ctx.Redirect(http.StatusFound, "/home")
		}
	} else {
		ctx.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "User not found",
		})
	}
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.Redirect(http.StatusFound, "/")
}
