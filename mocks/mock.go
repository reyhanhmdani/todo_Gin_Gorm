package mocks

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"todoGin/model/entity"
	"todoGin/repository"
)

// TodolistService adalah struktur yang merepresentasikan service Todolist.
type TodolistService struct {
	todolistRepository repository.TodoRepositoryMock
}

// NewTodolistService adalah fungsi constructor TodolistService.
func NewTodolistService(todolistRepository repository.TodoRepositoryMock) *TodolistService {
	return &TodolistService{todolistRepository: todolistRepository}
}

// GetAll adalah fungsi untuk mengambil semua todolist.
func (s *TodolistService) GetAll() ([]entity.Todolist, error) {
	todolist, err := s.todolistRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return todolist, nil
}

// GetByID adalah fungsi untuk mendapatkan data todolist berdasarkan ID.
func (s *TodolistService) GetByID(todoID int64) (*entity.Todolist, error) {
	// Cari todolist berdasarkan ID.
	todolist, err := s.todolistRepository.FindByID(todoID)
	if err != nil {
		return nil, err
	}

	return todolist, nil
}

// Create adalah fungsi untuk menambahkan todolist baru.
func (s *TodolistService) Create(title string) (*entity.Todolist, error) {
	if title == "" {
		return nil, errors.New("Title is required")
	}
	todolist := &entity.Todolist{Title: title}
	err := s.todolistRepository.Save(todolist)
	if err != nil {
		return nil, err
	}
	return todolist, nil
}
func (s *TodolistService) Update(todoID int64, updates map[string]interface{}) (*entity.Todolist, error) {
	// Ambil todolist yang akan diupdate.
	todolist, err := s.todolistRepository.FindByID(todoID)
	if err != nil {
		return nil, err
	}
	// Update todolist dengan field yang baru.
	for key, value := range updates {
		switch key {
		case "title":
			todolist.Title = value.(string)
		case "status":
			todolist.Status = value.(bool)
		default:
			return nil, errors.New("invalid field")
		}
	}
	// Simpan perubahan ke dalam repo
	err = s.todolistRepository.Upt(todolist)
	if err != nil {
		return nil, err
	}
	return todolist, nil
}

func (s *TodolistService) Delete(todoID int64) (int64, error) {
	// Cari todolist dengan ID yang diberikan.
	todolist, err := s.todolistRepository.FindByID(todoID)
	if err != nil {
		return 0, err
	}

	// Jika todolist tidak ditemukan, kembalikan error Not Found.
	if todolist == nil {
		return 0, errors.New("todolist not found")
	}

	// Hapus todolist dari repository.
	rowsAffected, err := s.todolistRepository.Del(todoID)
	if err != nil {
		return 0, errors.New("failed to delete todo item")
	}

	return rowsAffected, nil

}

// /////////////////////////////
type MockTodoRepository struct {
	mock.Mock
}

func (_m *MockTodoRepository) Save(todolist *entity.Todolist) error {
	args := _m.Called(todolist)
	return args.Error(0)

}
func (_m *MockTodoRepository) Del(todoID int64) (int64, error) {
	args := _m.Called(todoID)
	return args.Get(0).(int64), args.Error(1)
}

func (_m *MockTodoRepository) GetAll() ([]entity.Todolist, error) {
	args := _m.Called()
	return args.Get(0).([]entity.Todolist), args.Error(1)
}

// FindByID adalah implementasi mock untuk fungsi FindByID() di TodolistRepository.
func (_m *MockTodoRepository) FindByID(id int64) (*entity.Todolist, error) {
	args := _m.Called(id)
	return args.Get(0).(*entity.Todolist), args.Error(1)
}

// Update adalah implementasi mock untuk fungsi Update() di TodolistRepository.
func (_m *MockTodoRepository) Upt(todolist *entity.Todolist) error {
	args := _m.Called(todolist)
	return args.Error(0)
}
