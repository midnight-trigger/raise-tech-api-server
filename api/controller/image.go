package controller

import (
	"errors"
	"fmt"

	"github.com/labstack/echo"
	"github.com/midnight-trigger/raise-tech/api/definition"
	"github.com/midnight-trigger/raise-tech/api/domain"
	"github.com/midnight-trigger/raise-tech/logger"
	"github.com/midnight-trigger/raise-tech/third_party/jwt"
)

type Todo struct {
	Base
}

func (c *Todo) GetTodos(ctx echo.Context, claims *jwt.Claims) (response *Response) {
	defer func() {
		if e := recover(); e != nil {
			c.ServerErrorException(errors.New(""), fmt.Sprintf("%+v", e))
			logger.L.Error(c.ErrMessage)
			c.FormatResult(&c.Result, ctx)
		}
	}()

	params, err := definition.CreateGetTodosParam(ctx)
	if err != nil {
		c.ValidationException(err, err.Error())
		return c.FormatResult(&c.Result, ctx)
	}

	service := domain.GetNewTodoService()
	result := service.GetTodos(params, claims.UserId)
	return c.FormatResult(&result, ctx)
}
