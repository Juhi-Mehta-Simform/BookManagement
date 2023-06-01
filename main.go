package main

import (
	"Project/connection"
	"Project/models"
	"Project/router"
)

func main() {
	db := connection.GetConnection()
	//db.Migrator().DropTable(&models.Borrow{}, &models.Book{})
	db.AutoMigrate(&models.User{}, &models.Book{}, &models.Borrow{})
	defer connection.CloseConnection(db)
	router.Server.Run()
}
