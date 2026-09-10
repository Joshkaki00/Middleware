package main

import (
	"net"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// PreTrace runs in the Pre lane, before the router matches a route,
// so it fires even on requests that end up as a 404.
func PreTrace(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Response().Header().Set("X-Pre-Seen", "1")
		return next(c)
	}
}

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

// EnforceInside blocks any request in the group unless InsideTheBuilding
// already labeled it inside==true. Register AFTER InsideTheBuilding.
func EnforceInside(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		inside, _ := c.Get("inside").(bool)
		if !inside {
			return echo.NewHTTPError(http.StatusForbidden, "outside the building")
		}
		return next(c)
	}
}

func main() {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPDirect()

	e.Pre(PreTrace)

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

	internal := e.Group("/internal", EnforceInside)
	internal.GET("/admin", func(c *echo.Context) error {
		return c.String(http.StatusOK, "welcome, you are inside")
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}