package models

type Book struct {
	ID             int    `form:"id" json:"id" gorm:"primaryKey"`
	Title          string `form:"title" json:"title"`
	Author         string `form:"author" json:"author"`
	Genre          string `form:"genre" json:"genre"`
	Description    string `form:"description" json:"description"`
	ISBN           int    `form:"isbn" json:"isbn"`
	TotalQuantity  int    `form:"total_quantity" json:"total_quantity"`
	ActualQuantity int    `form:"actual_quantity" json:"actual_quantity"`
}
