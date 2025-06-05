package db

import (
	"fmt"
	"os"

	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB()*gorm.DB {
	host 		:= os.Getenv("DB_HOST")
	user 		:= os.Getenv("DB_USER")
	password	:= os.Getenv("DB_PASSWORD")
	dbname		:= os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v  sslmode=require", host, user, password, dbname)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }else {
		fmt.Println("Success connect to database!")
	}

	db.Debug().AutoMigrate(&models.Product{}, &models.User{})

	return db
}