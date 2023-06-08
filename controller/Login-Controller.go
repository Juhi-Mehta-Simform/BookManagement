package controller

import (
	"Project/connection"
	"Project/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/smtp"
	"strings"
)

var DB *gorm.DB

func init() {
	DB = connection.GetConnection()
}

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
	user.Gender = ctx.PostForm("gender")
	err := DB.Create(&user).Error
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
	DB.Model(&models.User{}).Where("email=?", email).Find(&user)
	if user.UserID != 0 {
		db := DB.Model(&models.User{}).Where("email=? AND password=?", email, password).Find(&user)
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
	db := DB.Model(&models.User{}).Where("email=?", email).Find(&user)
	if db.RowsAffected == 0 {
		ctx.HTML(http.StatusForbidden, "forgetPassword.html", gin.H{
			"error": "user not found",
		})
	} else {
		sendResetPasswordEmail(email)
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"error": "Password reset email send",
		})
	}
}

func sendResetPasswordEmail(email string) {
	auth := smtp.PlainAuth(
		"",
		"juhi.mehta.0604@gmail.com",
		"yczvyrzalemzefif",
		"smtp.gmail.com",
	)
	bytes, err := bcrypt.GenerateFromPassword([]byte(email), 0)
	msg := fmt.Sprintf("Subject: Reset Password\r\n"+
		"Please follow the below link to reset your password\n: http://localhost:8080/resetPassword?email=%s",
		string(bytes))
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"juhi.mehta.0604@gmail.com",
		[]string{email},
		[]byte(msg),
	)
	err = smtp.SendMail(
		"smtp.office365.com:587",
		auth,
		"juhi.mehta.0604@gmail.com",
		[]string{email},
		[]byte(msg),
	)
	if err != nil {
		log.Println(err)
	}

}
func LoadReset(ctx *gin.Context) {
	var users []models.User
	var user models.User
	email := ctx.Query("email")
	var userId int
	DB.Model(&models.User{}).Select("user_id, email").Find(&users)
	for i, _ := range users {
		err := bcrypt.CompareHashAndPassword([]byte(email), []byte(users[i].Email))
		if err == nil {
			userId = users[i].UserID
			break
		}
	}
	DB.Where("user_id=?", userId).Find(&user)
	ctx.HTML(http.StatusOK, "resetPassword.html", gin.H{
		"user": user,
	})
}

func ResetPassword(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	DB.Model(&models.User{}).Where("email=?", email).Update("password", password)
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"error": "password is reset.",
	})
}
