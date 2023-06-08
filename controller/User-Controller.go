package controller

import (
	"Project/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Home(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user models.User
	if session.Get("userID") != nil {
		id := session.Get("userID")
		err := DB.Model(&models.User{}).Where("user_id = ?", id).Find(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "home.html", gin.H{
				"user": user,
			})
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func ViewUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var users []models.User
	var user models.User
	var Roles []string
	if session.Get("userID") != nil {
		UserId := session.Get("userID")
		DB.Model(&models.User{}).Where("user_id = ?", UserId).Find(&user)
		err := DB.Model(&models.User{}).Order("user_id").Find(&users).Error
		DB.Model(&models.User{}).Distinct("role_name").Find(&Roles)
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
		DB.Model(&models.User{}).Where("user_id=?", userId).Update("role_name", "Librarian")
		ctx.Redirect(http.StatusFound, "/viewUser")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func RemoveLibrarian(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		userId := ctx.Param("user_id")
		DB.Model(&models.User{}).Where("user_id=?", userId).Update("role_name", "Member")
		ctx.Redirect(http.StatusFound, "/viewUser")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func MakeAdmin(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		userId := ctx.Param("user_id")
		DB.Model(&models.User{}).Where("user_id=?", userId).Update("role_name", "Admin")
		ctx.Redirect(http.StatusFound, "/viewUser")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

//func SearchUser(ctx *gin.Context) {
//	query := ctx.Query("query")
//	var users []models.User
//	DB.Model(&models.User{}).Where("name ILike ?", "%"+query+"%").Order("user_id").Find(&users)
//	ctx.JSON(http.StatusOK, users)
//}

func GetUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user models.User
	if session.Get("userID") != nil {
		userId := session.Get("userID")
		DB.Model(&models.User{}).Where("user_id = ?", userId).Find(&user)
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
		DB.Model(&models.User{}).Where("user_id = ?", userId).Find(&user)
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
		DB.Where("user_id=?", id).Find(&user)
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
			DB.Model(&models.User{}).Where("user_id=?", user.UserID).Updates(&user)
			ctx.Redirect(http.StatusFound, "/viewProfile")
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func SearchFilterUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var users []models.User
	var roleQuery string
	var search string
	if session.Get("userID") != nil {
		role := ctx.Query("role_name")
		query := ctx.Query("query")
		if len(role) != 0 && role != "" {
			tempRole := strings.Split(role, ",")
			roleQuery = "role_name IN ('" + strings.Join(tempRole, "','") + "')"
			search = " AND "
		}
		search += "name ILike " + "'%" + query + "%'"
		DB.Model(&models.User{}).Where(roleQuery + search).Order("user_id").Find(&users)
		ctx.JSON(http.StatusOK, users)
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}
