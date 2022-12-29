package main

import (
	"app/usecases"
	"fmt"
	"log"
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
	e.Static("/images", "images")

	api := e.Group("/api")

	images := api.Group("/images")
	images.POST("", upload)

	actuator := api.Group("/actuator")
	actuator.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func upload(c echo.Context) error {
	name := c.FormValue("name")

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	usecase := usecases.NewImageUsecase()

	err = usecase.InsertWatermark(src, fmt.Sprintf("images/%s", name))
	if err != nil {
		log.Fatal(err)
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf(`
		<p>File %s uploaded successfully with fields name=%s.</p>
		<img src="/images/%s" decode="async" alt="" />
		<a href="/">go back</a>
	`, file.Filename, name, name))
}
