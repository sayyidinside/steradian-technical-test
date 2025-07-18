package model

import (
	"steradian/interview/entity"
	"time"
)

type (
	OrderDetail struct {
		ID              int       `json:"id"`
		CarID           int       `json:"car_id"`
		CarName         string    `json:"car_name"`
		OrderDate       time.Time `json:"order_date"`
		PickupDate      time.Time `json:"pickup_date"`
		DropoffDate     time.Time `json:"dropoff_date"`
		PickupLocation  string    `json:"pickup_location"`
		DropoffLocation string    `json:"dropoff_location"`
	}

	OrderList struct {
		ID              int       `json:"id"`
		CarID           int       `json:"car_id"`
		CarName         string    `json:"car_name"`
		OrderDate       time.Time `json:"order_date"`
		PickupDate      time.Time `json:"pickup_date"`
		DropoffDate     time.Time `json:"dropoff_date"`
		PickupLocation  string    `json:"pickup_location"`
		DropoffLocation string    `json:"dropoff_location"`
	}

	OrderInput struct {
		CarID           int    `json:"car_id" form:"car_id" binding:"required"`
		OrderDate       string `json:"order_date" binding:"required"`
		PickupDate      string `json:"pickup_date" binding:"required"`
		DropoffDate     string `json:"dropoff_date" binding:"required"`
		PickupLocation  string `json:"pickup_location" binding:"required"`
		DropoffLocation string `json:"dropoff_location" binding:"required"`
	}
)

func (input *OrderInput) ToEntity() (*entity.Order, error) {
	orderDate, err := time.Parse("2006-01-02", input.OrderDate)
	if err != nil {
		return nil, err
	}

	pickupDate, err := time.Parse("2006-01-02", input.PickupDate)
	if err != nil {
		return nil, err
	}

	dropoffDate, err := time.Parse("2006-01-02", input.DropoffDate)
	if err != nil {
		return nil, err
	}

	return &entity.Order{
		CarID:           input.CarID,
		OrderDate:       orderDate,
		PickupDate:      pickupDate,
		DropoffDate:     dropoffDate,
		PickupLocation:  input.PickupLocation,
		DropoffLocation: input.DropoffLocation,
	}, nil
}
