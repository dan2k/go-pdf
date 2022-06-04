package orm

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
var Db *gorm.DB
var err error
func InitDB(){
	dsn := os.Getenv("DNS")
  	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err !=nil{
		panic(" Faild to connect database")
	}
	
	Db.AutoMigrate(&User{})
}