package entity

import (
	"time"
)

type Order struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CarID           int       `json:"car_id" gorm:"index"`
	OrderDate       time.Time `json:"order_date" gorm:"type:date;not null"`
	PickupDate      time.Time `json:"pickup_date" gorm:"type:date;not null"`
	DropoffDate     time.Time `json:"dropoff_date" gorm:"type:date;not null"`
	PickupLocation  string    `json:"pickup_location" gorm:"type:char(50);not null"`
	DropoffLocation string    `json:"dropoff_location" gorm:"type:char(50);not null"`
	Car             Car       `json:"car" gorm:"foreignKey:CarID"`
}

func (o *Order) TableNames() string {
	return "cars"
}
