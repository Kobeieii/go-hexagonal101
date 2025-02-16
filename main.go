package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kobeieii/adapters"
	"kobeieii/core"
)

func main() {
	app := fiber.New()

	db, err := gorm.Open(sqlite.Open("orders.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&core.Order{})

	orderRepository := adapters.NewGormOrderRepository(db)
	orderService := core.NewOrderService(orderRepository)
	orderHandler := adapters.NewHttpOrderHandler(orderService)

	app.Post("/", orderHandler.CreateOrder)

	app.Listen(":3000")
}