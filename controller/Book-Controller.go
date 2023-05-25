package controller

import (
	"Project/connection"
	"Project/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func LoadBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		fmt.Println(session.Get("userID"))
		ctx.HTML(http.StatusOK, "addBook.html", nil)
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func CreateBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		fmt.Println(session.Get("userID"))
		var book models.Book
		book.Title = ctx.PostForm("title")
		book.Author = ctx.PostForm("author")
		book.Genre = ctx.PostForm("genre")
		book.Description = ctx.PostForm("description")
		book.ISBN, _ = strconv.Atoi(ctx.PostForm("isbn"))
		book.TotalQuantity, _ = strconv.Atoi(ctx.PostForm("total_quantity"))
		book.ActualQuantity, _ = strconv.Atoi(ctx.PostForm("actual_quantity"))
		connection.GetConnection().Create(&book)
		ctx.Redirect(http.StatusFound, "/viewBook")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func ViewBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var book []models.Book
	if session.Get("userID") != nil {
		err := connection.GetConnection().Debug().Model(&models.Book{}).Order("id").Find(&book).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "viewBook.html", gin.H{
				"book": book,
			})
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func LoadUpdate(ctx *gin.Context) {
	var book models.Book
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		fmt.Println(session.Get("userID"))
		id := ctx.Param("id")
		connection.GetConnection().Where("id=?", id).Find(&book)
		ctx.HTML(http.StatusOK, "updateBook.html", gin.H{
			"id":   id,
			"book": book,
		})
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func UpdateBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var book models.Book
	if session.Get("userID") != nil {
		err := ctx.Bind(&book)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			connection.GetConnection().Model(&models.Book{}).Where("id=?", book.ID).Updates(&book)
			ctx.Redirect(http.StatusFound, "/viewBook")
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func LoadDelete(ctx *gin.Context) {
	var book models.Book
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		fmt.Println(session.Get("userID"))
		id := ctx.Param("id")
		connection.GetConnection().Where("id=?", id).Delete(&book)
		ctx.Redirect(http.StatusFound, "/viewBook")
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

//func DeleteBook(ctx *gin.Context) {
//	session := sessions.Default(ctx)
//	var book models.Book
//	if session.Get("userID") != nil {
//		err := ctx.Bind(&book)
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, "not ok")
//		} else {
//			connection.GetConnection().Model(&models.Book{}).Debug().Where("id=?", book.ID).Delete(&book)
//			ctx.JSON(http.StatusOK, "deleted")
//		}
//	} else {
//		ctx.Redirect(http.StatusMovedPermanently, "/")
//	}
//}
