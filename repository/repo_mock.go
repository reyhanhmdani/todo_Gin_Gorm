package repository

import (
	"todoGin/model/entity"
)

type TodoRepositoryMock interface {
	GetAll() ([]entity.Todolist, error)
	FindByID(todoID int64) (*entity.Todolist, error)
	Save(todolist *entity.Todolist) error
	Upt(todolist *entity.Todolist) error
	Del(todoID int64) (int64, error)
}
