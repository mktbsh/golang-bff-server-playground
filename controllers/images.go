package controllers

import (
	"app/usecases"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ImagesController struct {
	Controller
	usecase usecases.ImageUsecase
}

func NewImagesController() ImagesController {
	return ImagesController{
		usecase: usecases.NewImageUsecase(),
	}
}

func (con ImagesController) UploadImage(c echo.Context) error {
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

	err = con.usecase.InsertWatermark(src, fmt.Sprintf("__LOCAL__/images/%s", name))
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf(`
		<p>File %s uploaded successfully with fields name=%s.</p>
		<img src="/__LOCAL__/images/%s" decode="async" alt="" />
		<a href="/">go back</a>
	`, file.Filename, name, name))
}
