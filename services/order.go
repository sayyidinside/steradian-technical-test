package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"steradian/interview/model"
	"steradian/interview/repository"
	"strings"
	"time"
)

type OrderService interface {
	Get(ctx context.Context, order_id int) (*model.OrderDetail, error)
	GetAll(ctx context.Context) ([]model.OrderList, error)
	Create(ctx context.Context, input *model.OrderInput) error
	Update(ctx context.Context, order_id int, input *model.OrderInput) error
	Delete(ctx context.Context, order_id int) error
}

type orderService struct {
	carRepo   repository.CarRepo
	orderRepo repository.OrderRepo
}

func NewOrderService(
	carRepo repository.CarRepo,
	orderRepo repository.OrderRepo) OrderService {
	return &orderService{
		carRepo:   carRepo,
		orderRepo: orderRepo,
	}
}

func (s *orderService) Get(ctx context.Context, order_id int) (*model.OrderDetail, error) {
	if order_id == 0 {
		return nil, errors.New("id cannot be 0")
	}

	order := s.orderRepo.FindByID(ctx, order_id)
	if order == nil {
		return nil, errors.New("order not found")
	}

	orderModel := &model.OrderDetail{
		ID:              order.ID,
		CarID:           order.CarID,
		CarName:         order.Car.Name,
		OrderDate:       order.OrderDate,
		PickupDate:      order.PickupDate,
		DropoffDate:     order.DropoffDate,
		PickupLocation:  order.PickupLocation,
		DropoffLocation: order.DropoffLocation,
	}

	return orderModel, nil
}

func (s *orderService) GetAll(ctx context.Context) ([]model.OrderList, error) {
	orders := s.orderRepo.FindAll(ctx)
	if len(orders) == 0 {
		return nil, errors.New("order data not found")
	}

	orderModels := []model.OrderList{}
	for _, order := range orders {
		orderModels = append(orderModels, model.OrderList{
			ID:              order.ID,
			CarID:           order.CarID,
			CarName:         order.Car.Name,
			OrderDate:       order.OrderDate,
			PickupDate:      order.PickupDate,
			DropoffDate:     order.DropoffDate,
			PickupLocation:  order.PickupLocation,
			DropoffLocation: order.DropoffLocation,
		})
	}

	return orderModels, nil
}
func (s *orderService) Create(ctx context.Context, input *model.OrderInput) error {
	orderEntity, err := input.ToEntity()
	if err != nil {
		return err
	}

	car := s.carRepo.FindByID(ctx, input.CarID)
	if car == nil {
		return errors.New("invalid car_id input")
	}

	orderDate, _ := time.Parse("2006-01-02", input.OrderDate)
	pickUpDate, _ := time.Parse("2006-01-02", input.PickupDate)
	dropoffDate, _ := time.Parse("2006-01-02", input.DropoffDate)
	existingOrder, err := s.orderRepo.FindOverlappingOrdersByCarID(ctx, car.ID, pickUpDate, dropoffDate)
	if existingOrder.ID != 0 || err != nil {
		log.Println("#####################")
		return errors.New("car already booked")
	}

	// backdate validation
	backdateErr := []string{}
	if orderDate.After(pickUpDate) {
		backdateErr = append(backdateErr, "pickup date cannot be before order date")
	}
	if pickUpDate.After(dropoffDate) {
		backdateErr = append(backdateErr, "drop off cannot be before pickup date")
	}

	if len(backdateErr) != 0 {
		return fmt.Errorf("backdate valdiation error: %s", strings.Join(backdateErr[:], ", "))
	}

	if err := s.orderRepo.Create(ctx, orderEntity); err != nil {
		return err
	}

	return nil
}

func (s *orderService) Update(ctx context.Context, order_id int, input *model.OrderInput) error {
	order := s.orderRepo.FindByID(ctx, order_id)
	if order == nil {
		return errors.New("order data not found")
	}

	car := s.carRepo.FindByID(ctx, input.CarID)
	if car == nil {
		return errors.New("invalid car_id input")
	}

	orderEntity, err := input.ToEntity()
	if err != nil {
		return err
	}

	orderEntity.ID = order_id

	if err := s.orderRepo.UpdateByID(ctx, orderEntity); err != nil {
		return err
	}

	return nil
}

func (s *orderService) Delete(ctx context.Context, order_id int) error {
	order := s.orderRepo.FindByID(ctx, order_id)
	if order == nil {
		return errors.New("order data not found")
	}

	if err := s.orderRepo.DeleteByID(ctx, order); err != nil {
		return err
	}

	return nil
}
