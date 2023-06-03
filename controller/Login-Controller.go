package controller

import (
	"Project/connection"
	"Project/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/smtp"
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

func LoadForget(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "forgetPassword.html", nil)
}

func ForgetPassword(ctx *gin.Context) {
	var user models.User
	email := ctx.PostForm("email")
	err := connection.GetConnection().Model(&models.User{}).Where("email=?", email).Find(&user).Error
	if err != nil {
		ctx.String(http.StatusOK, "User not found")
	} else {
		sendResetPasswordEmail(email)
		ctx.String(http.StatusOK, "Password reset mail send")
	}
}

func sendResetPasswordEmail(email string) {
	auth := smtp.PlainAuth(
		"",
		"juhi.mehta.0604@gmail.com",
		"yczvyrzalemzefif",
		"smtp.gmail.com",
	)
	msg := fmt.Sprintf("Subject: Reset Password\r\n"+
		"Please follow the below link to reset your password\n: http://localhost:8080/resetPassword/%s",
		email)
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"juhi.mehta.0604@gmail.com",
		[]string{email},
		[]byte(msg),
	)
	if err != nil {
		fmt.Println(email)
		fmt.Println("error")
		fmt.Println(err)
	}
}

func LoadReset(ctx *gin.Context) {
	var user models.User
	email := ctx.Param("email")
	db := connection.GetConnection().Model(&models.User{}).Where("email=?", email).Find(&user)
	connection.CloseConnection(db)
	ctx.HTML(http.StatusOK, "resetPassword.html", gin.H{
		"user": user,
	})
}

func ResetPassword(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	db := connection.GetConnection().Debug().Model(&models.User{}).Where("email=?", email).Update("password", password)
	connection.CloseConnection(db)
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"error": "password is reset.",
	})
}
