package models

import "time"

type User struct {
	UserID   int    `form:"user_id" json:"user_id" gorm:"primaryKey"`
	Name     string `form:"name" json:"name"`
	Email    string `form:"email" json:"email" gorm:"unique"`
	Password string `form:"password" json:"password"`
	RoleName string `json:"role_name" form:"role_name"`
	Age      int    `json:"age" form:"age"`
	Gender   string `json:"gender" form:"gender"`
}

type Book struct {
	ID             int      `form:"id" json:"id" gorm:"primaryKey"`
	Title          string   `form:"title" json:"title"`
	Author         string   `form:"author" json:"author"`
	Genre          string   `form:"genre" json:"genre"`
	Description    string   `form:"description" json:"description"`
	ISBN           string   `form:"isbn" json:"isbn" gorm:"unique"`
	TotalQuantity  int      `form:"total_quantity" json:"total_quantity"`
	ActualQuantity int      `form:"actual_quantity" json:"actual_quantity"`
	Borrow         []Borrow `gorm:"references:ISBN;foreignKey:BookISBN;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Donate         []Donate `gorm:"references:ISBN;foreignKey:BookISBN;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Borrow struct {
	BorrowID    int       `gorm:"primaryKey" form:"borrow_id" json:"borrow_id"`
	MemberID    int       `form:"member_id" json:"member_id"`
	LibrarianID int       `form:"librarian_id" json:"librarian_id"`
	BookISBN    string    `form:"book_isbn" json:"book_isbn"`
	Status      string    `form:"status" json:"status"`
	IssueDate   time.Time `gorm:"type:date" form:"issue_date" json:"issue_date"`
	DueDate     time.Time `gorm:"type:date"  form:"due_date" json:"due_date"`
	Member      User      `gorm:"references:UserID;foreignKey:MemberID"`
	Librarian   User      `gorm:"references:UserID;foreignKey:LibrarianID"`
}

type Donate struct {
	DonateID   int       `gorm:"primaryKey" form:"donate_id" json:"donate_id"`
	MemberID   int       `form:"member_id" json:"member_id"`
	BookISBN   string    `form:"book_isbn" json:"book_isbn"`
	DonateDate time.Time `gorm:"type:date"  form:"donate_date" json:"donate_date"`
	Quantity   int       `form:"quantity" json:"quantity"`
	Member     User      `gorm:"references:UserID;foreignKey:MemberID"`
}

type BorrowBook struct {
	MemberID    int
	Name        string
	LibrarianID int
	Title       string
	Genre       string
	Author      string
	BookISBN    string
	Status      string
	IssueDate   time.Time
	DueDate     time.Time
}

type DonateBook struct {
	MemberID   int
	Name       string
	Title      string
	Genre      string
	Author     string
	BookISBN   string
	DonateDate time.Time
	Quantity   int
}
