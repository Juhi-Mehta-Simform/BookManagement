package main

import (
	"Project/connection"
	"Project/controller"
	"Project/models"
	"Project/router"
)

func main() {
	//db.Migrator().DropTable(&models.Donate{}, models.User{}, &models.Borrow{}, &models.Book{})
	controller.DB.AutoMigrate(&models.Donate{}, &models.User{}, &models.Book{}, &models.Borrow{})
	defer connection.CloseConnection(controller.DB)
	router.Server.Run()
}
