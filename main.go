package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	e := echo.New()
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//router_url := os.Getenv("THIS_SERVER_HTTPS_ADDRESS")

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "webhook router 2.0")
	})

	e.POST("/bot", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.POST("/setWebhook", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.POST("/deleteWebhook", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.Logger.Fatal(e.Start(":80"))
}
