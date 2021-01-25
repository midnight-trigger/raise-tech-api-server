package domain

import (
	"errors"

	"github.com/midnight-trigger/todo/api/definition"
	"github.com/midnight-trigger/todo/infra/mysql"
	"github.com/midnight-trigger/todo/logger"
)

type Todo struct {
	Base
	MTodos mysql.ITodos
}

func GetNewTodoService() *Todo {
	todo := new(Todo)
	todo.MUsers = mysql.GetNewUser()
	todo.MTodos = mysql.GetNewTodo()
	return todo
}

// Todo検索・一覧取得
func (s *Todo) GetTodos(params *definition.GetTodosParam, userId string) (r Result) {
	r.New()

	// 検索条件（クエリパラメータ）をもとにTodoを検索・一覧を取得
	todos, err := s.MTodos.FindByQuery(params, userId)
	if err != nil {
		r.ServerErrorException(errors.New(""), err.Error())
		logger.L.Error(err)
		return
	}

	// レスポンス作成
	var responses []*definition.GetTodosResponse
	for _, todo := range todos {
		response := new(definition.GetTodosResponse)
		s.SetStructOnSameField(todo, response)
		responses = append(responses, response)
	}
	r.Data = responses

	pagination := new(Pagination)
	s.SetStructOnSameField(params, pagination)

	// 検索条件に一致するTodo数を取得・レスポンスに含める
	pagination.Total, err = s.MTodos.GetTotalCount(params, userId)
	if err != nil {
		r.ServerErrorException(errors.New(""), err.Error())
		logger.L.Error(err)
		return
	}
	r.Pagination = pagination
	return
}
