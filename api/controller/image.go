package controller

import (
	"errors"
	"fmt"

	"github.com/labstack/echo"
	"github.com/midnight-trigger/raise-tech-api-server/api/domain"
	"github.com/midnight-trigger/raise-tech-api-server/logger"
)

type Image struct {
	Base
}

func (c *Image) HealthCheck(ctx echo.Context) (response *Response) {
	result := domain.Result{}
	result.New()
	return c.FormatResult(&result, ctx)
}

func (c *Image) PostImage(ctx echo.Context) (response *Response) {
	defer func() {
		if e := recover(); e != nil {
			c.ServerErrorException(errors.New(""), fmt.Sprintf("%+v", e))
			logger.L.Error(c.ErrMessage)
			c.FormatResult(&c.Result, ctx)
		}
	}()

	imageFileName, err := domain.GetImageFileName(ctx)
	if err != nil {
		c.ImageUploadException(err, err.Error())
		return c.FormatResult(&c.Result, ctx)
	}

	service := domain.GetNewImageService()
	result := service.PostImage(imageFileName)
	return c.FormatResult(&result, ctx)
}
