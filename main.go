package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	infrastructure "steradian/interview/Infrastructure"
	"steradian/interview/entity"
	"steradian/interview/model"
	"steradian/interview/repository"
	"steradian/interview/services"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := infrastructure.ConnectDB()
	if err != nil {
		log.Panicln(err)
	}

	db.AutoMigrate(entity.Car{})
	db.AutoMigrate(entity.Order{})

	// Bootstraping
	// Repo
	carRepo := repository.NewCarRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Service
	carService := services.NewCarService(carRepo, orderRepo)
	orderService := services.NewOrderService(carRepo, orderRepo)

	router := gin.Default()

	router.Static("/assets", "./assets")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Car Routes
	router.GET("/api/v1/car/:id", func(c *gin.Context) {
		id := c.Param("id")

		car_id, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid car id")
		}
		ctx := context.Background()
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		car, err := carService.Get(ctxTimeout, car_id)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
		}

		c.JSON(http.StatusOK, car)
	})
	router.GET("/api/v1/car", func(c *gin.Context) {
		ctx := context.Background()
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		cars, err := carService.GetAll(ctxTimeout)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
		}

		c.JSON(http.StatusOK, cars)
	})
	router.POST("/api/v1/car", func(c *gin.Context) {
		input := model.CarInput{}

		if err := c.ShouldBind(&input); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		file, _ := c.FormFile("image")
		input.ImageName = file.Filename

		c.SaveUploadedFile(file, "./assets/"+file.Filename)

		ctx := context.Background()
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := carService.Create(ctxTimeout, &input); err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		c.String(http.StatusCreated, "Car data successfully created")
	})
	router.PUT("/api/v1/car/:id", func(c *gin.Context) {
		id := c.Param("id")

		car_id, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid car id")
		}

		input := model.CarInput{}

		if err := c.ShouldBind(&input); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		file, err := c.FormFile("image")
		if err != nil {
			c.String(http.StatusBadRequest, "missing image file")
			return
		}
		input.ImageName = file.Filename

		c.SaveUploadedFile(file, "./assets/"+file.Filename)

		ctx := context.Background()
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := carService.Update(ctxTimeout, car_id, &input); err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("Car %s data successfully updated", input.Name))
	})
	router.DELETE("/api/v1/car/:id", func(c *gin.Context) {
		id := c.Param("id")

		car_id, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid car id")
		}

		ctx := context.Background()
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := carService.Delete(ctxTimeout, car_id); err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		c.String(http.StatusOK, "Car data successfully deleted")
	})

	// Order Routes
	router.GET("/api/v1/order/:id", func(c *gin.Context) {
		id := c.Param("id")

		order_id, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid car id")
		}

		ctx := context.Background()
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		car, err := orderService.Get(ctxTimeout, order_id)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
		}

		c.JSON(http.StatusOK, car)
	})
	router.GET("/api/v1/order", func(c *gin.Context) {
		ctx := context.Background()
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		orders, err := orderService.GetAll(ctxTimeout)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
		}

		c.JSON(http.StatusOK, orders)
	})
	router.POST("/api/v1/order", func(c *gin.Context) {
		input := model.OrderInput{}

		if err := c.ShouldBind(&input); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		ctx := context.Background()
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := orderService.Create(ctxTimeout, &input); err != nil {
			if err.Error() == "car already booked" || strings.Contains(err.Error(), "backdate valdiation error") {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		c.String(http.StatusCreated, "Order data successfully created")
	})
	router.PUT("/api/v1/order/:id", func(c *gin.Context) {
		id := c.Param("id")

		order_id, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid order id")
		}

		input := model.OrderInput{}

		if err := c.ShouldBind(&input); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		ctx := context.Background()
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := orderService.Update(ctxTimeout, order_id, &input); err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		c.String(http.StatusOK, "Order data succesfully updated")
	})
	router.DELETE("/api/v1/order/:id", func(c *gin.Context) {
		id := c.Param("id")

		order_id, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid order id")
		}

		ctx := context.Background()
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()

		if err := orderService.Delete(ctxTimeout, order_id); err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		c.String(http.StatusOK, "Order data successfully deleted")
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
