package models

type User struct {
	UserID   int    `form:"user_id" json:"user_id" gorm:"primaryKey"`
	Name     string `form:"name" json:"name"`
	Email    string `form:"email" json:"email" gorm:"unique"`
	Password string `form:"password" json:"password"`
	Borrow   []Borrow
}

type Borrow struct {
	BorrowID  int    `gorm:"primaryKey" form:"borrow_id" json:"borrow_id"`
	UserID    int    `form:"user_id" json:"user_id"`
	IssueDate string `gorm:"type:date" form:"issue_date" json:"issue_date"`
	DueDate   string `gorm:"type:date" form:"due_date" json:"due_date"`
}
