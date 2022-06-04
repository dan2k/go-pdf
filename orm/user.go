package orm

import "gorm.io/gorm"

type User struct {
	gorm.Model
	User     string
	Password string
	Fullname string
	Avatar   string
}