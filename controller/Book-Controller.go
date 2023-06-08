package controller

import (
	"Project/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func LoadBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		ctx.HTML(http.StatusOK, "addBook.html", nil)
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
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
		err := DB.Create(&book).Error
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				ctx.HTML(http.StatusBadRequest, "addBook.html", gin.H{
					"error": "Book already exists",
				})
			}
		} else {
			ctx.Redirect(http.StatusFound, "/viewBook")
		}
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
	data := ctx.Query("error")
	if session.Get("userID") != nil {
		userId := session.Get("userID")
		DB.Model(&models.User{}).Where("user_id = ?", userId).Find(&user)
		DB.Model(&models.Book{}).Distinct("genre").Find(&Genres)
		DB.Model(&models.Book{}).Distinct("author").Find(&Authors)
		err := DB.Model(&models.Book{}).Order("id").Find(&book).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "viewBook.html", gin.H{
				"book":    book,
				"user":    user,
				"authors": Authors,
				"genres":  Genres,
				"error":   data,
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
		DB.Where("id=?", id).Find(&book)
		ctx.HTML(http.StatusOK, "updateBook.html", gin.H{
			"id":   id,
			"book": book,
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
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
			DB.Model(&models.Book{}).Where("id=?", book.ID).Updates(&book)
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
		DB.Where("id=?", id).Delete(&book)
		ctx.Redirect(http.StatusFound, "/viewBook")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func LoadBorrow(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var book models.Book
	var borrow models.Borrow
	if session.Get("userID") != nil {
		borrow.MemberID = session.Get("userID").(int)
		borrow.BookISBN = ctx.Param("isbn")
		id := ctx.Param("id")
		DB.Where("id=?", id).Find(&book)
		db := DB.Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", borrow.MemberID, borrow.BookISBN).Find(&borrow)
		if db.RowsAffected != 0 {
			ctx.Redirect(http.StatusFound, "/viewBook?error=You already borrowed this book.")
		} else {
			ctx.HTML(http.StatusOK, "borrowBook.html", gin.H{
				"id":   id,
				"book": book,
			})
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func BorrowBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var borrow models.Borrow
	var book models.Book
	if session.Get("userID") != nil {
		borrow.MemberID = session.Get("userID").(int)
		borrow.LibrarianID, _ = strconv.Atoi(ctx.PostForm("librarian_id"))
		borrow.BookISBN = ctx.PostForm("book_isbn")
		borrow.IssueDate, _ = time.Parse("2006-01-02", ctx.PostForm("issue_date"))
		borrow.DueDate, _ = time.Parse("2006-01-02", ctx.PostForm("due_date"))
		borrow.Status = ctx.PostForm("status")
		db := DB.Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", borrow.MemberID, borrow.BookISBN).Find(&borrow)
		if db.RowsAffected == 0 {
			DB.Where("isbn=?", borrow.BookISBN).Find(&book)
			DB.Create(&borrow)
			newQuantity := book.ActualQuantity - 1
			DB.Model(&models.Book{}).Where("isbn=?", borrow.BookISBN).Update("actual_quantity", newQuantity)
			ctx.Redirect(http.StatusFound, "/userBorrow")
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func ViewBorrow(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var borrowList []models.BorrowBook
	var user models.User
	today := time.Now().Format("2006-01-02")
	if session.Get("userID") != nil {
		id := session.Get("userID").(int)
		DB.Model(&models.User{}).Where("user_id = ?", id).Find(&user)
		err := DB.Model(&models.Borrow{}).Select("b.title, b.genre, b.author, borrows.book_isbn, borrows.issue_date, borrows.due_date, u.name, borrows.status, borrows.librarian_id").Joins("JOIN books AS b ON borrows.book_isbn = b.isbn JOIN users AS u on borrows.member_id=u.user_id").Order("status, borrow_id").Find(&borrowList).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "allBorrow.html", gin.H{
				"borrow": borrowList,
				"user":   user,
				"today":  today,
			})
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func ViewRequest(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var borrowList []models.BorrowBook
	if session.Get("userID") != nil {
		err := DB.Model(&models.Borrow{}).Select("b.title, b.genre, b.author, borrows.book_isbn, borrows.issue_date, borrows.due_date, borrows.member_id, borrows.status, borrows.librarian_id").Joins("JOIN books AS b ON borrows.book_isbn = b.isbn").Order("status, borrow_id").Find(&borrowList).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "viewReturnRequest.html", gin.H{
				"borrow": borrowList,
			})
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func UserBorrow(ctx *gin.Context) {
	var borrowList []models.BorrowBook
	var user models.User
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		id := session.Get("userID").(int)
		DB.Model(&models.User{}).Where("user_id = ?", id).Find(&user)
		DB.Model(&models.Borrow{}).Where("member_id=?", id).Select("b.title, b.genre, b.author, borrows.book_isbn, borrows.issue_date, borrows.due_date, borrows.status, borrows.librarian_id, borrows.member_id").Joins("JOIN books AS b ON borrows.book_isbn = b.isbn").Order("status, borrow_id").Find(&borrowList)
		ctx.HTML(http.StatusOK, "viewBorrow.html", gin.H{
			"member_id": id,
			"borrow":    borrowList,
			"user":      user,
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func BorrowHistory(ctx *gin.Context) {
	var borrowList []models.BorrowBook
	today := time.Now().Format("2006-01-02")
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		id := ctx.Param("user_id")
		DB.Model(&models.Borrow{}).Where("member_id=?", id).Select("b.title, b.genre, b.author, borrows.book_isbn, borrows.issue_date, borrows.due_date, borrows.status, borrows.librarian_id, u.name").Joins("JOIN books AS b ON borrows.book_isbn = b.isbn JOIN users AS u on borrows.member_id=u.user_id").Order("status, borrow_id").Find(&borrowList)
		ctx.HTML(http.StatusOK, "viewBorrow.html", gin.H{
			"member_id": id,
			"borrow":    borrowList,
			"today":     today,
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func LoadReturn(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var borrow models.Borrow
	if session.Get("userID") != nil {
		id := session.Get("userID").(int)
		isbn := ctx.Param("book_isbn")
		DB.Where("member_id=? AND book_isbn=?", id, isbn).Find(&borrow)
		ctx.HTML(http.StatusOK, "returnBook.html", gin.H{
			"member_id": id,
			"borrow":    borrow,
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
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
		db := DB.Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", borrow.MemberID, borrow.BookISBN).Find(&borrow)
		if db.RowsAffected != 0 {
			newStatus := "Requested"
			DB.Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", borrow.MemberID, borrow.BookISBN).Update("status", newStatus)
			ctx.Redirect(http.StatusFound, "/userBorrow")
		} else {
			ctx.JSON(http.StatusForbidden, "not ok")
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func ReturnRequest(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var borrow models.Borrow
	var book models.Book
	if session.Get("userID") != nil {
		memberId := ctx.Param("member_id")
		LibrarianId := session.Get("userID")
		isbn := ctx.Param("book_isbn")
		DB.Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", memberId, isbn).Find(&borrow)
		DB.Model(&models.Borrow{}).Where("member_id=? AND book_isbn=?", memberId, isbn).Updates(map[string]interface{}{"status": "Returned", "librarian_id": LibrarianId})
		DB.Where("isbn=?", borrow.BookISBN).Find(&book)
		newQuantity := book.ActualQuantity + 1
		DB.Model(&models.Book{}).Where("isbn=?", borrow.BookISBN).Update("actual_quantity", newQuantity)
		ctx.Redirect(http.StatusFound, "/userBorrow")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}
func SearchBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		query := ctx.Query("query")
		var books []models.Book
		DB.Model(&models.Book{}).Where("title ILike ? OR isbn ILike ?", "%"+query+"%", "%"+query+"%").Order("id").Find(&books)
		ctx.JSON(http.StatusOK, books)
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func FilterBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var books []models.Book
	var authorQuery string
	var genreQuery string
	if session.Get("userID") != nil {
		genre := ctx.Query("genres")
		author := ctx.Query("authors")
		if len(author) != 0 && author != "" {
			tempAuthor := strings.Split(author, ",")
			authorQuery = "author IN ('" + strings.Join(tempAuthor, "','") + "')"
		}
		if len(genre) != 0 && genre != "" {
			if len(author) != 0 && author != "" {
				genreQuery = " AND "
			}
			tempGenre := strings.Split(genre, ",")
			genreQuery += "genre IN ('" + strings.Join(tempGenre, "','") + "')"
		}
		DB.Model(&models.Book{}).Where(authorQuery + genreQuery).Order("id").Find(&books)
		ctx.JSON(http.StatusOK, books)
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func LoadDonate(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user models.User
	if session.Get("userID") != nil {
		id := session.Get("userID")
		DB.Where("user_id=?", id).Find(&user)
		ctx.HTML(http.StatusOK, "donateBook.html", gin.H{
			"id": id,
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func DonateBook(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var donate models.Donate
	var book models.Book
	if session.Get("userID") != nil {
		donate.MemberID = session.Get("userID").(int)
		donate.BookISBN = ctx.PostForm("book_isbn")
		donate.DonateDate, _ = time.Parse("2006-01-02", ctx.PostForm("donate_date"))
		donate.Quantity, _ = strconv.Atoi(ctx.PostForm("quantity"))
		title := ctx.PostForm("title")
		author := ctx.PostForm("author")
		genre := ctx.PostForm("genre")
		description := ctx.PostForm("description")
		db := DB.Model(&models.Book{}).Where("isbn=?", donate.BookISBN).Find(&book)
		if db.RowsAffected > 0 {
			DB.Where("isbn=?", donate.BookISBN).Find(&book)
			DB.Create(&donate)
			newTotalQuantity := book.TotalQuantity + donate.Quantity
			newQuantity := book.ActualQuantity + donate.Quantity
			DB.Model(&models.Book{}).Where("isbn=?", donate.BookISBN).Updates(map[string]interface{}{"total_quantity": newTotalQuantity, "actual_quantity": newQuantity})
			ctx.Redirect(http.StatusFound, "/userDonate")
		} else {
			newBook := models.Book{Title: title, Author: author, Genre: genre, Description: description, ISBN: donate.BookISBN, TotalQuantity: donate.Quantity, ActualQuantity: donate.Quantity}
			DB.Create(&newBook)
			DB.Create(&donate)
			ctx.Redirect(http.StatusFound, "/userDonate")
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func UserDonate(ctx *gin.Context) {
	var user models.User
	var donateList []models.DonateBook
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		id := session.Get("userID").(int)
		DB.Model(&models.User{}).Where("user_id = ?", id).Find(&user)
		DB.Model(&models.Donate{}).Where("member_id=?", id).Select("b.title, b.genre, b.author, donates.book_isbn, donates.donate_date, donates.quantity").Joins("JOIN books AS b ON donates.book_isbn = b.isbn").Order("donate_id").Scan(&donateList)
		ctx.HTML(http.StatusOK, "viewDonate.html", gin.H{
			"member_id": id,
			"donate":    donateList,
			"user":      user,
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func ViewDonate(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var donateList []models.DonateBook
	var user models.User
	if session.Get("userID") != nil {
		id := session.Get("userID").(int)
		DB.Model(&models.User{}).Where("user_id = ?", id).Find(&user)
		err := DB.Model(&models.Donate{}).Select("b.title, b.genre, b.author, donates.book_isbn, donates.donate_date, donates.quantity, u.name").Joins("JOIN books AS b ON donates.book_isbn = b.isbn JOIN users AS u on donates.member_id=u.user_id").Order("donate_id").Scan(&donateList).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "not ok")
		} else {
			ctx.HTML(http.StatusOK, "viewDonate.html", gin.H{
				"donate": donateList,
				"user":   user,
			})
		}
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}

func DonateHistory(ctx *gin.Context) {
	var user models.User
	var donateList []models.DonateBook
	session := sessions.Default(ctx)
	if session.Get("userID") != nil {
		id := ctx.Param("user_id")
		DB.Model(&models.User{}).Where("user_id = ?", id).Find(&user)
		DB.Model(&models.Donate{}).Where("member_id=?", id).Select("b.title, b.genre, b.author, donates.book_isbn, donates.donate_date, donates.quantity, u.name").Joins("JOIN books AS b ON donates.book_isbn = b.isbn JOIN users AS u on donates.member_id=u.user_id").Order("donate_id").Scan(&donateList)
		ctx.HTML(http.StatusOK, "viewDonate.html", gin.H{
			"member_id": id,
			"donate":    donateList,
		})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	}
}
func SendReminder(email string) error {
	auth := smtp.PlainAuth(
		"",
		"juhi.mehta.0604@gmail.com",
		"yczvyrzalemzefif",
		"smtp.gmail.com",
	)
	msg := fmt.Sprintf("Subject: Reminder for overdue book\r\n" +
		"This is a reminder for overdue books. So please return it as soon as possible")
	to := []string{email}
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"juhi.mehta.0604@gmail.com",
		to,
		[]byte(msg),
	)
	err = smtp.SendMail(
		"smtp.office365.com:587",
		auth,
		"juhi.mehta.0604@gmail.com",
		to,
		[]byte(msg),
	)
	if err != nil {
		return err
	}
	return nil
}

func Reminder(ctx *gin.Context) {
	var email string
	id := ctx.Param("member_id")
	DB.Model(&models.User{}).Where("user_id", id).Select("email").Find(&email)
	SendReminder(email)
	ctx.Redirect(http.StatusFound, "/home")
}
