package controller

import (
	"Project/connection"
	"Project/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func LoadBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		ctx.HTML(http.StatusOK, "addBook.html", nil)
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func CreateBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		var book models.Book
		book.Title = ctx.PostForm("title")
		book.Author = ctx.PostForm("author")
		book.Genre = ctx.PostForm("genre")
		book.Description = ctx.PostForm("description")
		book.ISBN = ctx.PostForm("isbn")
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

func DeleteBook(ctx *gin.Context) {
	var book models.Book
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		id := ctx.Param("id")
		connection.GetConnection().Where("id=?", id).Delete(&book)
		ctx.Redirect(http.StatusFound, "/viewBook")
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func LoadBorrow(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var book models.Book
	if session.Get("userID") != nil {
		id := ctx.Param("id")
		connection.GetConnection().Where("id=?", id).Find(&book)
		if book.ActualQuantity <= 0 {
			ctx.Redirect(http.StatusFound, "/viewBook")
		} else {
			ctx.HTML(http.StatusOK, "borrowBook.html", gin.H{
				"id":   id,
				"book": book,
			})
		}
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func BorrowBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var borrow models.Borrow
	var book models.Book
	if session.Get("userID") != nil {
		connection.GetConnection().Where("isbn=?", borrow.BookISBN).Find(&book)
		fmt.Println(book.ActualQuantity)
		borrow.MemberID = session.Get("userID").(int)
		borrow.BookISBN = ctx.PostForm("book_isbn")
		borrow.IssueDate, _ = time.Parse("2006-01-02", ctx.PostForm("issue_date"))
		borrow.DueDate, _ = time.Parse("2006-01-02", ctx.PostForm("due_date"))
		borrow.Status = ctx.PostForm("status")
		db := connection.GetConnection().Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", borrow.MemberID, borrow.BookISBN).Find(&borrow)
		if db.RowsAffected == 0 {

			connection.GetConnection().Where("isbn=?", borrow.BookISBN).Find(&book)
			connection.GetConnection().Create(&borrow)
			newQuantity := book.ActualQuantity - 1
			connection.GetConnection().Model(&models.Book{}).Where("isbn=?", borrow.BookISBN).Update("actual_quantity", newQuantity)
			ctx.Redirect(http.StatusFound, "/viewBook")
		} else {
			ctx.JSON(http.StatusForbidden, "not ok")
		}
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func ViewBorrow(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var borrow []models.Borrow
	if session.Get("userID") != nil {
		err := connection.GetConnection().Debug().Model(&models.Borrow{}).Order("borrow_id").Find(&borrow).Error
		//datesWithoutTimestamps := make([]string, len(borrow))
		//for i, result := range borrow {
		//	datesWithoutTimestamps[i] = result.IssueDate.Format("2006-01-02")
		//	fmt.Println(datesWithoutTimestamps[i])
		//}
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "viewBorrow.html", gin.H{
				"borrow": borrow,
			})
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func UserBorrow(ctx *gin.Context) {
	var borrow []models.Borrow
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		id := session.Get("userID").(int)
		connection.GetConnection().Model(&models.Borrow{}).Where("member_id=?", id).Find(&borrow)
		ctx.HTML(http.StatusOK, "viewBorrow.html", gin.H{
			"member_id": id,
			"borrow":    borrow,
		})
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func LoadReturn(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var borrow models.Borrow
	if session.Get("userID") != nil {
		id := session.Get("userID").(int)
		connection.GetConnection().Where("member_id=?", id).Find(&borrow)
		fmt.Println(borrow.IssueDate)
		ctx.HTML(http.StatusOK, "returnBook.html", gin.H{
			"member_id": id,
			"borrow":    borrow,
		})
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func ReturnBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var borrow models.Borrow
	var book models.Book
	if session.Get("userID") != nil {
		connection.GetConnection().Where("isbn=?", borrow.BookISBN).Find(&book)
		borrow.MemberID = session.Get("userID").(int)
		borrow.BookISBN = ctx.PostForm("book_isbn")
		borrow.IssueDate, _ = time.Parse("2006-01-02", ctx.PostForm("issue_date"))
		borrow.DueDate, _ = time.Parse("2006-01-02", ctx.PostForm("due_date"))
		borrow.Status = ctx.PostForm("status")
		db := connection.GetConnection().Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", borrow.MemberID, borrow.BookISBN).Find(&borrow)
		if db.RowsAffected != 0 {
			connection.GetConnection().Where("isbn=?", borrow.BookISBN).Find(&book)
			newQuantity := book.ActualQuantity + 1
			connection.GetConnection().Model(&models.Book{}).Where("isbn=?", borrow.BookISBN).Update("actual_quantity", newQuantity)
			connection.GetConnection().Model(&models.Borrow{}).Where("isbn=?", borrow.BookISBN).Update("status", "returned")
			ctx.Redirect(http.StatusFound, "/viewBook")
		} else {
			ctx.JSON(http.StatusForbidden, "not ok")
		}
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}
