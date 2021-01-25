package mysql

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/midnight-trigger/todo/api/definition"
)

type Todos struct {
	Id        int64      `gorm:"type:bigint(20);primary_key;auto_increment"`
	UserId    string     `gorm:"type:char(36);not null"`
	Title     string     `gorm:"type:varchar(255);not null"`
	Body      string     `gorm:"type:text;not null"`
	Status    string     `gorm:"type:enum('todo','progress','finished');not null;default:'todo'"`
	CreatedAt time.Time  `gorm:"type:datetime;not null"`
	UpdatedAt time.Time  `gorm:"type:datetime;not null"`
	DeletedAt *time.Time `gorm:"type:datetime"`
}

//go:generate mockgen -source todos.go -destination mock_mysql/mock_todos.go
type ITodos interface {
	FindByQuery(params *definition.GetTodosParam, userId string) (todos []Todos, err error)
	GetTotalCount(params *definition.GetTodosParam, userId string) (total int, err error)
	FindById(id int64) (todo Todos, err error)
	Create(todo *Todos) (insertedTodo *Todos, err error)
	Update(oldParams Todos, updateParams map[string]interface{}) (err error)
	Delete(todo *Todos) (err error)
}

func GetNewTodo() *Todos {
	return &Todos{}
}

func (m *Todos) FindByQuery(params *definition.GetTodosParam, userId string) (todos []Todos, err error) {
	query := db.Table("todos").
		Select("*").
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL")
	query = buildFindByQuery(query, params)

	err = query.Offset(params.Offset).
		Limit(params.Limit).
		Order(fmt.Sprintf("created_at %s", params.Sort)).
		Scan(&todos).Error
	return
}

func buildFindByQuery(query *gorm.DB, params *definition.GetTodosParam) *gorm.DB {
	if params.Title != "" {
		query = query.Where("title LIKE ?", "%"+params.Title+"%")
	}
	if params.Body != "" {
		query = query.Where("body LIKE ?", "%"+params.Body+"%")
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	return query
}

func (m *Todos) GetTotalCount(params *definition.GetTodosParam, userId string) (count int, err error) {
	query := db.Table("todos").
		Select("*").
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL")
	query = buildFindByQuery(query, params)
	err = query.Count(&count).Error
	return
}

func (m *Todos) FindById(id int64) (todo Todos, err error) {
	err = db.Where("id = ?", id).First(&todo).Error
	return
}

func (m *Todos) Create(todo *Todos) (insertedTodo *Todos, err error) {
	err = db.Create(todo).Error
	insertedTodo = todo
	return
}

func (m *Todos) Update(oldParams Todos, updateParams map[string]interface{}) (err error) {
	err = db.Model(&oldParams).Updates(updateParams).Error
	return
}

func (m *Todos) Delete(todo *Todos) (err error) {
	err = db.Delete(todo).Error
	return
}
