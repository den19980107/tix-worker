package models

import "time"

type User struct {
	Id          int       `gorm:"column:id"`
	Username    string    `gorm:"column:username"`
	IdNumber    string    `gorm:"column:idNumber"`
	PhoneNumber string    `gorm:"column:phoneNumber"`
	CreatedAt   time.Time `gorm:"column:createdAt"`
}

func (User) TableName() string {
	return "public.User"
}
