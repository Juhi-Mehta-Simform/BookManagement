package main

import (
	"Project/connection"
	"Project/models"
	"Project/router"
)

func main() {
	db := connection.GetConnection()
	db.AutoMigrate(&models.User{}, &models.Book{})
	defer connection.CloseConnection(db)

	router.Server.Run()
}
