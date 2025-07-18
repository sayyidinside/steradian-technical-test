package services

import (
	"context"
	"errors"
	"steradian/interview/model"
	"steradian/interview/repository"
)

type CarService interface {
	Get(ctx context.Context, car_id int) (*model.CarDetail, error)
	GetAll(ctx context.Context) ([]model.CarList, error)
	Create(ctx context.Context, input *model.CarInput) error
	Update(ctx context.Context, car_id int, input *model.CarInput) error
	Delete(ctx context.Context, car_id int) error
}

type carService struct {
	carRepo  repository.CarRepo
	oderRepo repository.OrderRepo
}

func NewCarService(
	carRepo repository.CarRepo,
	orderRepo repository.OrderRepo) CarService {
	return &carService{
		carRepo:  carRepo,
		oderRepo: orderRepo,
	}
}

func (s *carService) Get(ctx context.Context, car_id int) (*model.CarDetail, error) {
	if car_id == 0 {
		return nil, errors.New("id cannot be 0")
	}

	car := s.carRepo.FindByID(ctx, car_id)
	if car == nil {
		return nil, errors.New("car not found")
	}

	orderList := []model.OrderList{}
	for _, v := range car.Orders {
		orderList = append(orderList, model.OrderList{
			ID:        v.ID,
			CarID:     v.CarID,
			CarName:   car.Name,
			OrderDate: v.OrderDate,
		})
	}

	carModel := &model.CarDetail{
		ID:        car.ID,
		Name:      car.Name,
		DayRate:   car.DayRate,
		MonthRate: car.MonthRate,
		Image:     car.Image,
		Orders:    orderList,
	}

	return carModel, nil
}

func (s *carService) GetAll(ctx context.Context) ([]model.CarList, error) {
	cars := s.carRepo.FindAll(ctx)
	if len(cars) == 0 {
		return nil, errors.New("car data not found")
	}

	carModels := []model.CarList{}
	for _, car := range cars {
		carModels = append(carModels, model.CarList{
			ID:    car.ID,
			Name:  car.Name,
			Image: car.Image,
		})
	}

	return carModels, nil
}
func (s *carService) Create(ctx context.Context, input *model.CarInput) error {
	carEntity, err := input.ToEntity()
	if err != nil {
		return err
	}

	if err := s.carRepo.Create(ctx, carEntity); err != nil {
		return err
	}

	return nil
}

func (s *carService) Update(ctx context.Context, car_id int, input *model.CarInput) error {
	car := s.carRepo.FindByID(ctx, car_id)
	if car == nil {
		return errors.New("car data not found")
	}

	carEntity, err := input.ToEntity()
	if err != nil {
		return err
	}

	carEntity.ID = car_id

	if err := s.carRepo.UpdateByID(ctx, carEntity); err != nil {
		return err
	}

	return nil
}

func (s *carService) Delete(ctx context.Context, car_id int) error {
	car := s.carRepo.FindByID(ctx, car_id)
	if car == nil {
		return errors.New("car data not found")
	}

	if err := s.carRepo.DeleteByID(ctx, car); err != nil {
		return err
	}

	return nil
}
