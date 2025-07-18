package repository

import (
	"context"
	"steradian/interview/entity"

	"gorm.io/gorm"
)

type CarRepo interface {
	FindByID(ctx context.Context, car_id int) *entity.Car
	FindAll(ctx context.Context) (cars []entity.Car)
	Create(ctx context.Context, car *entity.Car) error
	UpdateByID(ctx context.Context, car *entity.Car) error
	DeleteByID(ctx context.Context, car *entity.Car) error
}

type carRepository struct {
	*gorm.DB
}

func NewCarRepository(db *gorm.DB) CarRepo {
	return &carRepository{DB: db}
}

func (r *carRepository) FindByID(ctx context.Context, car_id int) (car *entity.Car) {
	r.DB.WithContext(ctx).Model(entity.Car{}).Where("id = ?", car_id).Preload("Orders").Find(&car)
	return car
}

func (r *carRepository) FindAll(ctx context.Context) (cars []entity.Car) {
	r.DB.WithContext(ctx).Model(&entity.Car{}).Find(&cars)
	return cars
}

func (r *carRepository) Create(ctx context.Context, car *entity.Car) error {
	return r.DB.WithContext(ctx).Model(&entity.Car{}).Create(&car).Error
}

func (r *carRepository) UpdateByID(ctx context.Context, car *entity.Car) error {
	return r.DB.WithContext(ctx).Model(&entity.Car{}).Where("id = ?", car.ID).Updates(&car).Error
}

func (r *carRepository) DeleteByID(ctx context.Context, car *entity.Car) error {
	return r.DB.WithContext(ctx).Model(&entity.Car{}).Where("id = ?", car.ID).Delete(car).Error
}
