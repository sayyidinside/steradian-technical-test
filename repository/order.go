package repository

import (
	"context"
	"steradian/interview/entity"
	"time"

	"gorm.io/gorm"
)

type OrderRepo interface {
	FindByID(ctx context.Context, order_id int) *entity.Order
	FindAll(ctx context.Context) (orders []entity.Order)
	Create(ctx context.Context, order *entity.Order) error
	UpdateByID(ctx context.Context, order *entity.Order) error
	DeleteByID(ctx context.Context, order *entity.Order) error
	FindOverlappingOrdersByCarID(ctx context.Context, car_id int, pickupDate, dropoffDate time.Time) (*entity.Order, error)
}

type orderRepository struct {
	*gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepo {
	return &orderRepository{DB: db}
}

func (r *orderRepository) FindByID(ctx context.Context, order_id int) (order *entity.Order) {
	r.DB.WithContext(ctx).Model(entity.Order{}).Where("id = ?", order_id).Preload("Car", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Find(&order)
	return order
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order) error {
	return r.DB.WithContext(ctx).Model(&entity.Order{}).Create(&order).Error
}

func (r *orderRepository) UpdateByID(ctx context.Context, order *entity.Order) error {
	return r.DB.WithContext(ctx).Model(&entity.Order{}).Where("id = ?", order.ID).Updates(order).Error
}

func (r *orderRepository) DeleteByID(ctx context.Context, order *entity.Order) error {
	return r.DB.WithContext(ctx).Model(&entity.Order{}).Where("id = ?", order.ID).Delete(order).Error
}

func (r *orderRepository) FindAll(ctx context.Context) (orders []entity.Order) {
	r.DB.WithContext(ctx).Model(&entity.Order{}).Preload("Car", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Find(&orders)
	return orders
}

func (r *orderRepository) FindOverlappingOrdersByCarID(ctx context.Context, car_id int, pickupDate, dropoffDate time.Time) (*entity.Order, error) {
	var order entity.Order
	err := r.DB.WithContext(ctx).
		Model(&entity.Order{}).
		Where("car_id = ?", car_id).
		Where("pickup_date < ? AND dropoff_date > ?", dropoffDate, pickupDate).
		Preload("Car", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Limit(1).
		Find(&order).Error
	return &order, err
}
