package controller

import (
	"Project/connection"
	"Project/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Home(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user models.User
	if session.Get("userID") != nil {
		id := session.Get("userID")
		err := connection.GetConnection().Debug().Model(&models.User{}).Where("user_id = ?", id).Find(&user).Error
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
		db := connection.GetConnection().Create(&book)
		defer connection.CloseConnection(db)
		ctx.Redirect(http.StatusFound, "/viewBook")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func ViewBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var book []models.Book
	var user models.User
	var Authors []string
	var Genres []string
	if session.Get("userID") != nil {
		userId := session.Get("userID")
		db := connection.GetConnection().Debug().Model(&models.User{}).Where("user_id = ?", userId).Find(&user)
		defer connection.CloseConnection(db)
		err := connection.GetConnection().Debug().Model(&models.Book{}).Order("id").Find(&book).Error
		db = connection.GetConnection().Model(&models.Book{}).Distinct("genre").Find(&Genres)
		db = connection.GetConnection().Model(&models.Book{}).Distinct("author").Find(&Authors)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "viewBook.html", gin.H{
				"book":    book,
				"user":    user,
				"authors": Authors,
				"genres":  Genres,
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
		db := connection.GetConnection().Where("id=?", id).Find(&book)
		defer connection.CloseConnection(db)
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
			db := connection.GetConnection().Model(&models.Book{}).Where("id=?", book.ID).Updates(&book)
			defer connection.CloseConnection(db)
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
		db := connection.GetConnection().Where("id=?", id).Delete(&book)
		defer connection.CloseConnection(db)
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
		db := connection.GetConnection().Where("id=?", id).Find(&book)
		defer connection.CloseConnection(db)
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
		db := connection.GetConnection().Where("isbn=?", borrow.BookISBN).Find(&book)
		defer connection.CloseConnection(db)
		fmt.Println(book.ActualQuantity)
		borrow.MemberID = session.Get("userID").(int)
		borrow.BookISBN = ctx.PostForm("book_isbn")
		borrow.IssueDate, _ = time.Parse("2006-01-02", ctx.PostForm("issue_date"))
		borrow.DueDate, _ = time.Parse("2006-01-02", ctx.PostForm("due_date"))
		borrow.Status = ctx.PostForm("status")
		db = connection.GetConnection().Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", borrow.MemberID, borrow.BookISBN).Find(&borrow)
		if db.RowsAffected == 0 {
			db = connection.GetConnection().Where("isbn=?", borrow.BookISBN).Find(&book)
			db = connection.GetConnection().Create(&borrow)
			newQuantity := book.ActualQuantity - 1
			db = connection.GetConnection().Model(&models.Book{}).Where("isbn=?", borrow.BookISBN).Update("actual_quantity", newQuantity)
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
	var user models.User
	if session.Get("userID") != nil {
		id := session.Get("userID").(int)
		db := connection.GetConnection().Debug().Model(&models.User{}).Where("user_id = ?", id).Find(&user)
		defer connection.CloseConnection(db)
		fmt.Println(user.RoleName)
		err := connection.GetConnection().Debug().Model(&models.Borrow{}).Order("status, borrow_id").Find(&borrow).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "viewBorrow.html", gin.H{
				"borrow": borrow,
				"user":   user,
			})
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func UserBorrow(ctx *gin.Context) {
	var borrow []models.Borrow
	var user models.User
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		id := session.Get("userID").(int)
		db := connection.GetConnection().Debug().Model(&models.User{}).Where("user_id = ?", id).Find(&user)
		defer connection.CloseConnection(db)
		fmt.Println(user.RoleName)
		db = connection.GetConnection().Model(&models.Borrow{}).Where("member_id=?", id).Order("status, borrow_id").Find(&borrow)
		ctx.HTML(http.StatusOK, "viewBorrow.html", gin.H{
			"member_id": id,
			"borrow":    borrow,
			"user":      user,
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
		isbn := ctx.Param("book_isbn")
		db := connection.GetConnection().Debug().Where("member_id=? AND book_isbn=?", id, isbn).Find(&borrow)
		fmt.Println(isbn)
		defer connection.CloseConnection(db)
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
	if session.Get("userID") != nil {
		borrow.MemberID = session.Get("userID").(int)
		borrow.BookISBN = ctx.PostForm("book_isbn")
		borrow.DueDate, _ = time.Parse("2006-01-02", ctx.PostForm("due_date"))
		borrow.Status = ctx.PostForm("status")
		db := connection.GetConnection().Debug().Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", borrow.MemberID, borrow.BookISBN).Find(&borrow)
		defer connection.CloseConnection(db)
		if db.RowsAffected != 0 {
			newStatus := "requested"
			db = connection.GetConnection().Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", borrow.MemberID, borrow.BookISBN).Update("status", newStatus)
			ctx.Redirect(http.StatusFound, "/viewBook")
		} else {
			ctx.JSON(http.StatusForbidden, "not ok")
		}
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func ReturnRequest(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var borrow models.Borrow
	var book models.Book
	if session.Get("userID") != nil {
		memberId := ctx.Param("member_id")
		isbn := ctx.Param("book_isbn")
		db := connection.GetConnection().Debug().Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", memberId, isbn).Find(&borrow)
		defer connection.CloseConnection(db)
		db = connection.GetConnection().Debug().Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", memberId, isbn).Update("status", "Returned")
		db = connection.GetConnection().Where("isbn=?", borrow.BookISBN).Find(&book)
		newQuantity := book.ActualQuantity + 1
		db = connection.GetConnection().Model(&models.Book{}).Where("isbn=?", borrow.BookISBN).Update("actual_quantity", newQuantity)
		ctx.Redirect(http.StatusFound, "/viewBook")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}
func SearchBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		query := ctx.Query("query")
		var books []models.Book
		db := connection.GetConnection().Model(&models.Book{}).Where("title ILike ?  OR author ILike ?  OR genre ILike ?  OR isbn ILike ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").Order("id").Find(&books)
		defer connection.CloseConnection(db)
		ctx.JSON(http.StatusOK, books)
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func FilterBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var db *gorm.DB
	var books []models.Book
	if session.Get("userID") != nil {
		//genre := ctx.Query("genres")
		author := ctx.Query("authors")
		fmt.Println(author)
		if len(author) != 0 && author != "" {
			db = connection.GetConnection().Debug().Model(&models.Book{}).Where("author IN (?)", strings.Split(author, ",")).Find(&books)
		} else {
			db = connection.GetConnection().Find(&books)
		}
		fmt.Println(books)
		connection.CloseConnection(db)
		ctx.JSON(http.StatusOK, books)
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}

}
