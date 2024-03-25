package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Order struct {
	Id           int       `gorm:"column:id"`
	From         Station   `gorm:"column:from"`
	To           Station   `gorm:"column:to"`
	CreatorId    int       `gorm:"column:creatorId"`
	CreatedAt    time.Time `gorm:"column:createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt"`
	DepartureDay time.Time `gorm:"column:departureDay"`
	StartTime    string    `gorm:"column:startTime"`
	EndTime      string    `gorm:"column:endTime"`
	ExecDay      string    `gorm:"column:execDay"`
	Creator      User      `gorm:"foreignKey:CreatorId"`
	Captcha      string    `gorm:"column:captcha"`
	JsessionId   string    `gorm:"column:jsessionId"`
}

func (Order) TableName() string {
	return "Order"
}

func (o Order) GetStartTime() string {
	hour, _ := strconv.Atoi(strings.Split(o.StartTime, ":")[0])
	min, _ := strconv.Atoi(strings.Split(o.StartTime, ":")[1])

	minStr := ""
	if min < 30 {
		minStr = "00"
	} else {
		minStr = "30"
	}

	meridiem := "A"
	if hour > 11 {
		hour -= 12
		meridiem = "P"
	}

	return fmt.Sprintf("%d%s%s", hour, minStr, meridiem)
}
