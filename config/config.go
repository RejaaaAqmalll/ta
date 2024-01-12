package config

import (
	"fmt"
	"os"
	"ta-kasir/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
func ConnectDatabase() *gorm.DB {
	godotenv.Load()
	sqlString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_HOST"),
	os.Getenv("DB_PORT"),
	os.Getenv("DB_NAME"),
)

	db, err := gorm.Open(mysql.Open(sqlString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(
		&model.User{},
	)

	if err != nil {
		panic(err)
	}

DB = db
return DB
}