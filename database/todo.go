package database

import (
	"errors"
	"gorm.io/gorm"
	"todoGin/model/entity"
	"todoGin/repository"
)

// adaptop pattern
type TodoRepository struct {
	DB *gorm.DB
}

func NewTodoRepository(dbClient *gorm.DB) repository.TodoRepository {
	return &TodoRepository{
		DB: dbClient,
	}
}

func (t TodoRepository) GetAll() ([]entity.Todolist, error) {
	var todos []entity.Todolist

	result := t.DB.Find(&todos)
	return todos, result.Error
}

func (t TodoRepository) GetByID(todoID int64) (*entity.Todolist, error) {
	var todo entity.Todolist
	result := t.DB.Where("id = ?", todoID).First(&todo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &todo, result.Error
}

func (t TodoRepository) Create(title string) (*entity.Todolist, error) {
	todo := entity.Todolist{
		Title: title,
	}
	result := t.DB.Create(&todo)
	return &todo, result.Error
}

func (t TodoRepository) Update(todoID int64, updates map[string]interface{}) (*entity.Todolist, error) {
	var todo entity.Todolist
	result := t.DB.Model(&todo).Where("id = ?", todoID).Updates(updates)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &todo, result.Error
}

func (t TodoRepository) Delete(todoID int64) (int64, error) {
	todo := entity.Todolist{ID: todoID}
	result := t.DB.Delete(&todo)
	return result.RowsAffected, result.Error
	//if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	//	return false, nil
	//}
	//fmt.Printf("Id: %d", todoID)
	//fmt.Println(result.Error)
	//return true, result.Error
}
