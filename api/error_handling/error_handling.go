package error_handling

import (
	"fmt"
	"net/http"
)

type ErrorHandling struct {
	Code          int
	ErrMessage    string
	RawErrMessage string
	ErrStack      error
}

func (e *ErrorHandling) CognitoErrorFoundException(stackErr error, errMessage string) {
	e.Code = http.StatusBadRequest
	e.ErrMessage = errMessage
	e.setSlackErrorInfo(stackErr, "")
}

func (e *ErrorHandling) setSlackErrorInfo(stackErr error, rawErrMessage string) {
	e.ErrStack = stackErr
	e.RawErrMessage = rawErrMessage
}

func (e *ErrorHandling) ServerErrorException(stackErr error, rawErrMessage string) {
	e.Code = http.StatusInternalServerError
	e.ErrMessage = "サーバーエラー"
	e.setSlackErrorInfo(stackErr, rawErrMessage)
}

func (e *ErrorHandling) ValidationException(stackErr error, errMessage string) {
	e.Code = http.StatusBadRequest
	e.ErrMessage = errMessage
	e.setSlackErrorInfo(stackErr, "")
}

func (e *ErrorHandling) TodoNotFoundException(stackErr error) {
	e.Code = http.StatusNotFound
	e.ErrMessage = "対象Todoが見つかりません"
	e.setSlackErrorInfo(stackErr, "")
}

func (e *ErrorHandling) UserIsNotOwnerException(stackErr error) {
	e.Code = http.StatusBadRequest
	e.ErrMessage = "必要な権限がありません"
	e.setSlackErrorInfo(stackErr, "")
}

func GetValidationErrorMessage(field string, tag string, params string) (message string) {
	switch tag {
	case "required":
		switch field {
		case "Username":
			message = "ユーザ名は必須です"
		}
	case "email":
		switch field {
		case "Email":
			message = "メールアドレスの形式が正しくありません"
		}
	case "gte":
		switch field {
		case "Password":
			message = "パスワードは8文字以上で入力して下さい"
		}
	case "lte":
		switch field {
		case "Username":
			message = fmt.Sprintf("ユーザ名は%s文字以内で入力して下さい", params)
		}
	case "oneof":
		switch field {
		case "Sort":
			message = "ソートはDESC, ASCのいずれかで入力して下さい"
		}
	}
	return
}
