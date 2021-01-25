package definition

import (
	"strconv"

	"github.com/labstack/echo"
	"github.com/midnight-trigger/todo/logger"
	"github.com/pkg/errors"
)

type PutTodoParam struct {
	TodoId int64 `validate:"required"`
}

type PutTodoRequestBody struct {
	Title string `json:"title" validate:"required,lte=255"`
	Body  string `json:"body" validate:"required"`
}

type PutTodoResponse struct {
	Id        int64  `json:"id"`
	UserId    string `json:"user_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func CreatePutTodoRequestBodyAndParam(ctx echo.Context) (param *PutTodoParam, body *PutTodoRequestBody, err error) {
	todoId := ctx.Param("todoId")
	param = new(PutTodoParam)
	param.TodoId, _ = strconv.ParseInt(todoId, 10, 64)
	if message, ok := Validator(param); !ok {
		err = errors.New(message)
		logger.L.Error(message)
		return
	}

	body = new(PutTodoRequestBody)
	if err = ctx.Bind(body); err != nil {
		logger.L.Error(err)
		return
	}
	if message, ok := Validator(body); !ok {
		err = errors.New(message)
		logger.L.Error(message)
	}
	return
}
