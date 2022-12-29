package main

import (
	"app/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Static("/", "public")
	e.Static("/__LOCAL__", "__LOCAL__")

	api := e.Group("/api")

	images := api.Group("/images")
	{
		imagesController := controllers.NewImagesController()
		images.POST("", imagesController.UploadImage)
	}

	actuator := api.Group("/actuator")
	{
		actuator.GET("/health", func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})
	}

	e.Logger.Fatal(e.Start(":1323"))
}
