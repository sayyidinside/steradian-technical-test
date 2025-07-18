package entity

import "github.com/shopspring/decimal"

type Car struct {
	ID        int             `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string          `json:"name" gorm:"type:char(50);not null"`
	DayRate   decimal.Decimal `json:"day_rate" gorm:"type:double;not null"`
	MonthRate decimal.Decimal `json:"month_rate" gorm:"type:double;not null"`
	Image     string          `json:"image" gorm:"type:char(255);not null"`
	Orders    []Order         `json:"orders" gorm:"foreignKey:CarID"`
}

func (C *Car) TableNames() string {
	return "cars"
}
