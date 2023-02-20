package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"todoGin/mocks"
	//"github.com/stretchr/testify/mock"
	"testing"
	"todoGin/model/entity"
)

func TestGetAllSuccess(t *testing.T) {
	// Buat mock TodolistRepository.
	mockRepo := new(mocks.MockTodoRepository)

	// Buat slice dari Todolist.
	expectedTodolists := []entity.Todolist{
		{ID: 1, Title: "Belajar Golang", Status: true},
		{ID: 2, Title: "Belajar Unit Testing", Status: false},
		{ID: 3, Title: "Belajar Mocking", Status: false},
	}

	// Atur return value untuk mock GetAll().
	mockRepo.On("GetAll").Return(expectedTodolists, nil)

	// Buat service dengan mock repository.
	service := mocks.NewTodolistService(mockRepo)

	// Panggil GetAll() dari service.
	todolists, err := service.GetAll()

	// Periksa apakah hasil panggilan GetAll() sesuai dengan ekspektasi.
	assert.NoError(t, err)
	assert.Equal(t, expectedTodolists, todolists)

	// Periksa apakah mock GetAll() dipanggil dengan benar.
	mockRepo.AssertCalled(t, "GetAll")
}

// Failed case:
func TestGetAllInternalServerError(t *testing.T) {
	// Buat mock TodolistRepository.
	mockRepo := new(mocks.MockTodoRepository)

	mockTodolists := []entity.Todolist{
		entity.Todolist{ID: 1, Title: "Todo 1", Status: false},
		entity.Todolist{ID: 2, Title: "Todo 2", Status: true},
	}

	// Atur return value untuk mock FindAll().
	mockRepo.On("GetAll").Return(mockTodolists, nil)

	// Buat service dengan mock repository.
	service := mocks.NewTodolistService(mockRepo)

	// Panggil GetAll() dari service.
	todolist, err := service.GetAll()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, mockTodolists, todolist)

	// Ensure that the mocks expectations were met
	mockRepo.AssertExpectations(t)
}

func TestCreateSuccess(t *testing.T) {
	// Buat mock TodolistRepository.
	mockRepo := new(mocks.MockTodoRepository)

	// Buat Todolist baru.
	newTodolist := &entity.Todolist{Title: "Belajar Golang"}

	// Atur return value untuk mock Save().
	mockRepo.On("Save", newTodolist).Return(nil)

	// Buat service dengan mock repository.
	s := mocks.NewTodolistService(mockRepo)

	// Panggil Create() dari service.
	todolist, err := s.Create(newTodolist.Title)

	// Periksa apakah hasil panggilan Create() sesuai dengan ekspektasi.
	assert.NoError(t, err)
	assert.Equal(t, newTodolist.Title, todolist.Title)

	// Periksa apakah mock Save() dipanggil dengan benar.
	mockRepo.AssertCalled(t, "Save", newTodolist)
}
func TestCreateTodolistInvalid(t *testing.T) {
	//membuat mock object TodolistRepository
	mockRepo := new(mocks.MockTodoRepository)

	//mengembalikan error Invalid jika title kosong
	mockRepo.On("Save", "").Return(nil, errors.New("Title is required"))

	//melakukan testing
	//service := TodolistService{repository: mockRepo}
	service := mocks.NewTodolistService(mockRepo)
	result, err := service.Create("")

	//memastikan error yang dikembalikan sama dengan error Invalid
	assert.Equal(t, errors.New("Title is required"), err)
	assert.Nil(t, result)

	//memastikan method Create pada mock object tidak dipanggil
	mockRepo.AssertNotCalled(t, "Save")
}
func TestCreateInternalServerError(t *testing.T) {
	// Buat mock TodolistRepository.
	mockRepo := new(mocks.MockTodoRepository)

	// Atur return value untuk mock Create().
	expectedError := errors.New("database error")
	// Atur mock Save() untuk selalu mengembalikan error yang tidak di inginkan.
	mockRepo.On("Save", mock.AnythingOfType("*entity.Todolist")).Return(expectedError)

	// Buat service dengan mock repository.
	s := mocks.NewTodolistService(mockRepo)

	// Panggil Create() dari service.
	_, err := s.Create("Belajar Golang")

	// Periksa apakah error yang dihasilkan sesuai dengan ekspektasi.
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())

	// Periksa apakah mock Create() dipanggil dengan benar.
	mockRepo.AssertCalled(t, "Save", mock.AnythingOfType("*entity.Todolist"))
}

func TestUpdateInternalServerError(t *testing.T) {
	mockRepo := new(mocks.MockTodoRepository)
	expectedErr := errors.New("INTERNAL SERVER ERROR")

	// Buat Todolist yang akan diperbarui.
	todolist := &entity.Todolist{ID: 1, Title: "Belajar Golang", Status: false}

	// Atur return value untuk mock FindByID().
	mockRepo.On("FindByID", todolist.ID).Return(todolist, expectedErr)

	mockRepo.On("Upt", int64(1), map[string]interface{}{"Title": "Updated Task", "Description": "This task has been updated"}).Return(nil, expectedErr)
	service := mocks.NewTodolistService(mockRepo)

	todo, err := service.Update(1, map[string]interface{}{"Title": "Updated Task", "Description": "This task has been updated"})
	assert.NotNil(t, err)
	assert.EqualError(t, errors.New("INTERNAL SERVER ERROR"), err.Error())
	assert.Nil(t, todo)
}

// Update
func TestUpdateSuccess(t *testing.T) {
	// Buat mock TodolistRepository.
	mockRepo := new(mocks.MockTodoRepository)

	// Buat Todolist yang akan diperbarui.
	todolist := &entity.Todolist{ID: 1, Title: "Belajar Golang", Status: false}

	// Atur return value untuk mock FindByID().
	mockRepo.On("FindByID", todolist.ID).Return(todolist, nil)

	// Atur return value untuk mock Update().
	mockRepo.On("Upt", todolist).Return(nil)

	// Buat service dengan mock repository.
	s := mocks.NewTodolistService(mockRepo)

	// Update Todolist.
	updates := map[string]interface{}{
		"title":  "Belajar Golang framework",
		"status": true,
	}
	updatedTodolist, err := s.Update(todolist.ID, updates)

	// Periksa apakah hasil panggilan Update() sesuai dengan ekspektasi.
	assert.NoError(t, err)
	assert.Equal(t, "Belajar Golang framework", updatedTodolist.Title)
	//assert.Equal

}

func TestUpdateNotFound(t *testing.T) {
	// Buat mock TodolistRepository.
	mockRepo := new(mocks.MockTodoRepository)

	todolistupt := &entity.Todolist{ID: 1, Title: "Belajar Golang", Status: false}

	// Atur return value untuk mock FindByID().
	mockRepo.On("FindByID", int64(111)).Return(todolistupt, errors.New("not found"))

	// Buat service dengan mock repository.
	service := mocks.NewTodolistService(mockRepo)

	// Panggil Update() dari service.
	todolist, err := service.Update(111, map[string]interface{}{"title": "Belajar Golang"})

	// Periksa apakah error yang dihasilkan sesuai dengan ekspektasi.
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "not found")

	// Periksa apakah todolist yang dihasilkan nil.
	assert.Nil(t, todolist)

	// Periksa apakah mock FindByID() dipanggil dengan benar.
	mockRepo.AssertCalled(t, "FindByID", int64(111))
}

func TestGetByID(t *testing.T) {
	// Buat mock TodolistRepository.
	mockRepo := new(mocks.MockTodoRepository)

	// Atur return value untuk mock FindByID().
	mockTodo := &entity.Todolist{
		ID:     1,
		Title:  "Belajar Golang",
		Status: true,
	}
	mockRepo.On("FindByID", int64(1)).Return(mockTodo, nil)

	// Buat service dengan mock repository.
	service := mocks.NewTodolistService(mockRepo)

	// Panggil GetByID() dari service.
	todolist, err := service.GetByID(1)

	// Periksa apakah todolist yang dihasilkan sesuai dengan ekspektasi.
	assert.NoError(t, err)
	assert.NotNil(t, todolist)
	assert.Equal(t, mockTodo, todolist)

	// Periksa apakah mock FindByID() dipanggil dengan benar.
	mockRepo.AssertCalled(t, "FindByID", int64(1))
}

func TestGetByIDNotFound(t *testing.T) {
	// Buat mock TodolistRepository.
	mockRepo := new(mocks.MockTodoRepository)

	mockTodo := &entity.Todolist{
		ID:     1,
		Title:  "Belajar Golang",
		Status: true,
	}

	// Atur return value untuk mock FindByID().
	mockRepo.On("FindByID", int64(99)).Return(mockTodo, errors.New("not found"))

	// Buat service dengan mock repository.
	service := mocks.NewTodolistService(mockRepo)

	// Panggil GetByID() dari service.
	todolist, err := service.GetByID(99)

	// Periksa apakah error yang dihasilkan sesuai dengan ekspektasi.
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "not found")

	// Periksa apakah todolist yang dihasilkan nil.
	assert.Nil(t, todolist)

	// Periksa apakah mock FindByID() dipanggil dengan benar.
	mockRepo.AssertCalled(t, "FindByID", int64(99))
}

func TestGetByIDInternalServerError(t *testing.T) {
	// Buat mock TodolistRepository.
	mockRepo := new(mocks.MockTodoRepository)

	data := &entity.Todolist{ID: 1, Title: "Belajar Golang", Status: false}

	// Atur return value untuk mock FindByID().
	mockRepo.On("FindByID", int64(1)).Return(data, errors.New("internal server error"))

	// Buat service dengan mock repository.
	service := mocks.NewTodolistService(mockRepo)

	// Panggil GetByID() dari service.
	todolist, err := service.GetByID(1)

	// Periksa apakah error yang dihasilkan sesuai dengan ekspektasi.
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "internal server error")

	// Periksa apakah todolist yang dihasilkan nil.
	assert.Nil(t, todolist)

	// Periksa apakah mock FindByID() dipanggil dengan benar.
	mockRepo.AssertCalled(t, "FindByID", int64(1))
}

func TestDeleteSuccess(t *testing.T) {

	// membuat mock object TodoRepo
	mockRepo := new(mocks.MockTodoRepository)
	// Buat Todolist yang akan diperbarui.
	todolist := &entity.Todolist{ID: 1, Title: "Belajar Golang", Status: false}

	// Atur return value untuk mock FindByID().
	mockRepo.On("FindByID", todolist.ID).Return(todolist, nil)
	mockRepo.On("Del", int64(1)).Return(int64(1), nil)

	// membuat service dan memasukkan mock object TodoRepo
	service := mocks.NewTodolistService(mockRepo)

	// memanggil fungsi Delete
	rowsAffected, err := service.Delete(int64(1))

	// memastikan mock object dipanggil dengan parameter yang benar
	mockRepo.AssertCalled(t, "Del", int64(1))

	// memastikan hasil pengembalian dari fungsi Delete adalah benar
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)
}

func TestDeleteInternalServerError(t *testing.T) {
	// membuat mock object TodoRepo
	mockRepo := new(mocks.MockTodoRepository)
	// Buat Todolist yang akan diperbarui.
	todolist := &entity.Todolist{ID: 1, Title: "Belajar Golang", Status: false}

	// Atur return value untuk mock FindByID().
	mockRepo.On("FindByID", todolist.ID).Return(todolist, nil)
	mockRepo.On("Del", int64(1)).Return(int64(0), errors.New("INTERNAL SERVER ERROR"))

	// membuat service dan memasukkan mock object TodoRepo
	service := mocks.NewTodolistService(mockRepo)

	// memanggil fungsi Delete
	rowsAffected, err := service.Delete(int64(1))

	// memastikan mock object dipanggil dengan parameter yang benar
	mockRepo.AssertCalled(t, "Del", int64(1))

	// memastikan hasil pengembalian dari fungsi Delete adalah benar
	assert.Error(t, err)
	assert.Equal(t, int64(0), rowsAffected)
}

func TestDeleteNotFound(t *testing.T) {
	// membuat mock object TodoRepo
	mockRepo := new(mocks.MockTodoRepository)
	// Buat Todolist yang akan diperbarui.
	todolist := &entity.Todolist{ID: 1, Title: "Belajar Golang", Status: false}

	// Atur return value untuk mock FindByID().
	mockRepo.On("FindByID", int64(11)).Return(todolist, errors.New("not found"))

	mockRepo.On("Del").Return(nil)

	// membuat service dan memasukkan mock object TodoRepo
	service := mocks.NewTodolistService(mockRepo)

	_, err := service.Delete(int64(11))

	// Periksa apakah error yang dihasilkan sesuai dengan ekspektasi.
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "not found")

	//assert.Nil(t, todo)

	// Periksa apakah mock FindByID() dipanggil dengan benar.
	mockRepo.AssertCalled(t, "FindByID", int64(11))
}
