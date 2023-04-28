package main

import (
	"net/http"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Setup Main Server
	echoMainServer := echo.New()
	echoMainServer.HideBanner = true
	echoMainServer.Use(middleware.Logger())
	echoMainServer.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, balancer!")
	})
	echoMainServer.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})
	// Create Prometheus server and Middleware
	echoPrometheus := echo.New()
	echoPrometheus.HideBanner = true
	prom := prometheus.NewPrometheus("echo", nil)
	// Scrape metrics from Main Server
	echoMainServer.Use(prom.HandlerFunc)
	// Setup metrics endpoint at another server
	prom.SetMetricsPath(echoPrometheus)
	go func() { echoPrometheus.Logger.Fatal(echoPrometheus.Start(":8081")) }()
	echoMainServer.Logger.Fatal(echoMainServer.Start(":8080"))
}
