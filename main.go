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

func InsideTheBuilding(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ipStr := c.RealIP()
		ip := net.ParseIP(ipStr)
		inside := false
		via := "unparseable"
		if ip != nil {
			inside = ip.IsLoopback() || ip.IsPrivate()
			via = "realip+private/loopback"
		}
		c.Set("inside", inside)
		c.Set("via", via)
		c.Response().Header().Set("X-Inside-The-Building", map[bool]string{true: "1", false: "0"}[inside])
		return next(c)
	}
}

func main() {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPDirect()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(ServerHeader)
	e.Use(InsideTheBuilding)

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/whoami", func(c *echo.Context) error {
		inside, _ := c.Get("inside").(bool)
		via, _ := c.Get("via").(string)
		return c.JSON(http.StatusOK, map[string]any{
			"ip":     c.RealIP(),
			"inside": inside,
			"via":    via,
		})
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}