package main

import (
	"net"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "New-Echo-Server")
		return next(c)
	}
}

func main() {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPDirect()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(ServerHeader)

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/whoami", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"ip": c.RealIP()})
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}