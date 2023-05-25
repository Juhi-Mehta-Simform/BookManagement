package connection

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetConnection() *gorm.DB {
	dbURI := fmt.Sprintf("host=localhost user=postgres dbname=Project sslmode=disable password=root port=5432")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURI,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	errCheck(err)
	return db
}

func CloseConnection(db *gorm.DB) {
	dbs, err := db.DB()
	errCheck(err)
	err = dbs.Close()
	errCheck(err)
}
