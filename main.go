package main

import (
	"context"
	"os"

	"github.com/tifye/remote-led/core"
	"github.com/tifye/remote-led/view"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var (
	port string
)

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "assets",
		Browse: false,
	}))
	e.Use(middleware.Logger())

	e.GET("/", index)
	e.POST("/actions/on", turnOn)
	e.POST("/actions/off", turnOff)

	e.Logger.Fatal(e.Start(":" + port))
}

func index(ctx echo.Context) error {
	isOn, err := core.IsOn(context.Background())
	if err != nil {
		log.Errorf("failed to get isOn: %v", err)
		return err
	}

	comp := view.Index(isOn)
	return comp.Render(context.Background(), ctx.Response().Writer)
}

func turnOff(ctx echo.Context) error {
	err := core.Fill(context.Background(), core.NewRGB(0, 0, 0))
	if err != nil {
		log.Errorf("failed to fill: %v", err)
		return err
	}

	comp := view.Switch(false)
	return comp.Render(context.Background(), ctx.Response().Writer)
}

func turnOn(ctx echo.Context) error {
	err := core.Fill(context.Background(), core.NewRGB(255, 255, 255))
	if err != nil {
		log.Errorf("failed to fill: %v", err)
		return err
	}

	comp := view.Switch(true)
	return comp.Render(context.Background(), ctx.Response().Writer)
}
