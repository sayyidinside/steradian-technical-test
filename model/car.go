package model

import (
	"steradian/interview/entity"

	"github.com/shopspring/decimal"
)

type (
	CarDetail struct {
		ID        int             `json:"id"`
		Name      string          `json:"name"`
		DayRate   decimal.Decimal `json:"day_rate"`
		MonthRate decimal.Decimal `json:"month_rate"`
		Image     string          `json:"image"`
		Orders    []OrderList     `json:"orders"`
	}

	CarList struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	CarInput struct {
		Name      string `json:"name" form:"name" binding:"required"`
		DayRate   string `json:"day_rate" form:"day_rate" binding:"required"`
		MonthRate string `json:"month_rate" form:"month_rate" binding:"required"`
		ImageName string `json:"image_name" form:"image_name"`
	}
)

func (input *CarInput) ToEntity() (*entity.Car, error) {
	dayRate, err := decimal.NewFromString(input.DayRate)
	if err != nil {
		return nil, err
	}

	monthRate, err := decimal.NewFromString(input.MonthRate)
	if err != nil {
		return nil, err
	}

	return &entity.Car{
		Name:      input.Name,
		DayRate:   dayRate,
		MonthRate: monthRate,
		Image:     input.ImageName,
	}, nil
}
