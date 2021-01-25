package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/midnight-trigger/todo/api/domain"
	"github.com/midnight-trigger/todo/third_party/slack"
)

type Base struct {
	domain.Result
}

type Response struct {
	Meta       *Meta              `json:"meta"`
	Data       interface{}        `json:"data"`
	Pagination *domain.Pagination `json:"pagination,omitempty"`
}

type Meta struct {
	Code         int    `json:"code"`
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`
}

func (b *Base) FormatResult(r *domain.Result, ctx echo.Context) *Response {
	response := new(Response)

	switch r.Code {
	case http.StatusOK: //200
		meta := &Meta{
			Code: r.Code,
		}
		response.Meta = meta
		response.Data = r.Data
		response.Pagination = r.Pagination

	case http.StatusBadRequest: //400
		fallthrough
	case http.StatusUnauthorized: //401
		fallthrough
	case http.StatusForbidden: //403
		fallthrough
	case http.StatusNotFound: //404
		fallthrough
	case http.StatusConflict: //409
		fallthrough
	case http.StatusInternalServerError: //500
		meta := &Meta{
			Code:         r.Code,
			ErrorType:    http.StatusText(r.Code),
			ErrorMessage: r.ErrMessage,
		}
		response.Meta = meta

		body, _ := ioutil.ReadAll(ctx.Request().Body)
		slack.SlackSend(ctx.Request().RequestURI, ctx.Request().Method, r.Code, ctx.Request().Header, string(body), meta.ErrorType, meta.ErrorMessage, fmt.Sprintf("%+v", r.ErrStack), r.RawErrMessage)
	}

	return response
}
