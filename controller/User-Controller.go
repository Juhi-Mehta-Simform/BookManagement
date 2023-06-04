package controller

import (
	"Project/connection"
	"Project/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func ViewUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var users []models.User
	var user models.User
	var Roles []string
	if session.Get("userID") != nil {
		UserId := session.Get("userID")
		db := connection.GetConnection().Debug().Model(&models.User{}).Where("user_id = ?", UserId).Find(&user)
		defer connection.CloseConnection(db)
		err := connection.GetConnection().Model(&models.User{}).Order("user_id").Find(&users).Error
		db = connection.GetConnection().Model(&models.User{}).Distinct("role_name").Find(&Roles)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "viewUser.html", gin.H{
				"users": users,
				"user":  user,
				"roles": Roles,
			})
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func MakeLibrarian(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		userId := ctx.Param("user_id")
		connection.GetConnection().Debug().Model(&models.User{}).Where("user_id=?", userId).Update("role_name", "Librarian")
		ctx.Redirect(http.StatusFound, "/viewUser")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func MakeAdmin(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		userId := ctx.Param("user_id")
		connection.GetConnection().Debug().Model(&models.User{}).Where("user_id=?", userId).Update("role_name", "Admin")
		ctx.Redirect(http.StatusFound, "/viewUser")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func SearchUser(ctx *gin.Context) {
	query := ctx.Query("query")
	fmt.Println(query)
	var users []models.User
	db := connection.GetConnection().Model(&models.User{}).Where("name ILike ?  OR role_name ILike ?", "%"+query+"%", "%"+query+"%").Order("user_id").Find(&users)
	defer connection.CloseConnection(db)
	fmt.Println(users)
	ctx.JSON(http.StatusOK, users)
}

func GetUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user models.User
	if session.Get("userID") != nil {
		userId := session.Get("userID")
		db := connection.GetConnection().Debug().Model(&models.User{}).Where("user_id = ?", userId).Find(&user)
		defer connection.CloseConnection(db)
		ctx.JSON(http.StatusOK, user)
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func UserProfile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user models.User
	if session.Get("userID") != nil {
		userId := session.Get("userID")
		db := connection.GetConnection().Debug().Model(&models.User{}).Where("user_id = ?", userId).Find(&user)
		defer connection.CloseConnection(db)
		ctx.HTML(http.StatusOK, "viewProfile.html", gin.H{
			"user": user,
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func LoadProfile(ctx *gin.Context) {
	var user models.User
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		id := ctx.Param("user_id")
		db := connection.GetConnection().Where("user_id=?", id).Find(&user)
		defer connection.CloseConnection(db)
		ctx.HTML(http.StatusOK, "updateProfile.html", gin.H{
			"user_id": id,
			"user":    user,
		})
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func UpdateProfile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user models.User
	if session.Get("userID") != nil {
		err := ctx.Bind(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			fmt.Println(user)
			db := connection.GetConnection().Debug().Model(&models.User{}).Where("user_id=?", user.UserID).Updates(&user)
			defer connection.CloseConnection(db)
			ctx.Redirect(http.StatusFound, "/viewProfile")
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func FilterUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var db *gorm.DB
	var users []models.User
	if session.Get("userID") != nil {
		role := ctx.Query("role_name")
		if len(role) != 0 && role != "" {
			db = connection.GetConnection().Debug().Model(&models.User{}).Where("role_name IN (?)", strings.Split(role, ",")).Order("user_id").Find(&users)
		} else {
			db = connection.GetConnection().Order("user_id").Find(&users)
		}
		connection.CloseConnection(db)
		ctx.JSON(http.StatusOK, users)
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}
