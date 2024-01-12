package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/tifye/remote-led/core"
	"github.com/tifye/remote-led/view"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	regPerSec = 20
)

var (
	port         string
	ledServerUrl string
	ledService   *core.LedService
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Warnf("failed to load .env file")
	}

	port = os.Getenv("PORT")
	if port == "" {
		log.Printf("PORT env is not set, using default port 9000")
		port = "9000"
	}

	ledServerUrl = os.Getenv("LED_SERVER_URL")
	if ledServerUrl == "" {
		log.Fatal("LED_SERVER_URL env is not set")
	}
}

func main() {
	ledService = core.NewLedService(ledServerUrl)

	e := echo.New()

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(regPerSec)))
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
	isOn, err := ledService.IsOn(context.Background())
	if err != nil {
		log.Errorf("failed to get isOn: %v", err)
		return err
	}

	comp := view.Index(isOn)
	return comp.Render(context.Background(), ctx.Response().Writer)
}

func turnOff(ctx echo.Context) error {
	err := ledService.Fill(context.Background(), core.NewRGB(0, 0, 0))
	if err != nil {
		log.Errorf("failed to fill: %v", err)
		return err
	}

	comp := view.Switch(false)
	return comp.Render(context.Background(), ctx.Response().Writer)
}

func turnOn(ctx echo.Context) error {
	err := ledService.Fill(context.Background(), core.NewRGB(255, 255, 255))
	if err != nil {
		log.Errorf("failed to fill: %v", err)
		return err
	}

	comp := view.Switch(true)
	return comp.Render(context.Background(), ctx.Response().Writer)
}
