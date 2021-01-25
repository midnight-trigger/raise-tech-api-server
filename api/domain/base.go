package domain

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/midnight-trigger/todo/api/error_handling"
	"github.com/midnight-trigger/todo/configs"
	"github.com/midnight-trigger/todo/logger"
)

type Base struct {
}

type Result struct {
	Data       interface{}
	Pagination *Pagination
	error_handling.ErrorHandling
}

type Pagination struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (r *Result) New() {
	r.Code = http.StatusOK
}

func (b *Base) SetStructOnSameField(fromStruct interface{}, toStruct interface{}) {
	from := reflect.Indirect(reflect.ValueOf(fromStruct))
	to := reflect.Indirect(reflect.ValueOf(toStruct))

	for i := 0; i < from.NumField(); i++ {
		fromFieldValue := from.Field(i).Interface()
		if t, ok := to.Type().FieldByName(from.Type().Field(i).Name); ok {
			switch v := fromFieldValue.(type) {
			case int64:
				to.Field(t.Index[0]).SetInt(v)
			case int:
				to.Field(t.Index[0]).SetInt(int64(v))
			case string:
				to.Field(t.Index[0]).SetString(v)
			case time.Time:
				to.Field(t.Index[0]).SetString(v.Format("2006-01-02 15:04:05"))
			}
		}
	}
	return
}

func (b *Base) TestInit(t *testing.T) (ctrl *gomock.Controller) {
	ctrl = gomock.NewController(t)
	defer ctrl.Finish()
	configs.Init("test")
	logger.Init("test")
	return
}
