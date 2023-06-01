package controller

import (
	"Project/connection"
	"Project/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ViewUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var users []models.User
	if session.Get("userID") != nil {
		err := connection.GetConnection().Model(&models.User{}).Order("user_id").Find(&users).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "viewUser.html", gin.H{
				"users": users,
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

func SearchUser(ctx *gin.Context) {
	query := ctx.Query("query")
	fmt.Println(query)
	var users []models.User
	db := connection.GetConnection().Model(&models.User{}).Where("name ILike ?  OR role_name ILike ?", "%"+query+"%", "%"+query+"%").Order("user_id").Find(&users)
	defer connection.CloseConnection(db)
	fmt.Println(users)
	ctx.JSON(http.StatusOK, users)
}
