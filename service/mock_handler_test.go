package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"todoGin/model/entity"
	"todoGin/model/request"
	"todoGin/model/respErr"
	"todoGin/repository"
)

func TestTodolist(t *testing.T) {
	t.Run("testGetAllSuccess", testGetAllSuccess)
	t.Run("testTodolistInternalServerError", testTodolistInternalServerError)
	t.Run("testEmptyTodo", testEmptyTodo)
	t.Run("testCreateSuccess", testCreateSuccess)
	t.Run("testCreateInternalServerError", testCreateInternalServerError)
	t.Run("testCreateInvalid", testCreateInvalid)
}

func testGetAllSuccess(t *testing.T) {
	todoRepo := repository.NewMockTodoRepository(t)

	// Buat slice dari Todolist.
	expectedTodolists := []entity.Todolist{
		{ID: 1, Title: "Belajar Golang", Status: true},
		{ID: 2, Title: "Belajar Unit Testing", Status: false},
		{ID: 3, Title: "Belajar Mocking", Status: false},
	}

	todoRepo.On("GetAll").Return(expectedTodolists, nil)

	todoSvc := NewTodoService(todoRepo)
	endpoint := "/manage-todos"

	r := gin.New()
	r.GET(endpoint, todoSvc.TodolistHandlerGetAll)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	respBody, err := io.ReadAll(w.Body)
	require.NoError(t, err)
	//
	var result request.TodoResponseToGetAll
	err = json.Unmarshal(respBody, &result)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedTodolists, result.Todos)

}

func testTodolistInternalServerError(t *testing.T) {
	todoRepo := repository.NewMockTodoRepository(t)

	// Buat slice dari Todolist.
	expectedInternalServerTodolists := []entity.Todolist{
		{ID: 1, Title: "Belajar Golang", Status: true},
		{ID: 2, Title: "Belajar Unit Testing", Status: false},
		{ID: 3, Title: "Belajar Mocking", Status: false},
	}
	expectedErr := errors.New("Internal Server Error")
	todoRepo.On("GetAll").Return(expectedInternalServerTodolists, expectedErr)

	todoSvc := NewTodoService(todoRepo)
	endpoint := "/manage-todos"

	r := gin.New()
	r.GET(endpoint, todoSvc.TodolistHandlerGetAll)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	respBody, err := io.ReadAll(w.Body)
	require.NoError(t, err)

	var errResp respErr.ErrorResponse
	err = json.Unmarshal(respBody, &errResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedErr.Error(), errResp.Message)
}

func testEmptyTodo(t *testing.T) {
	todoRepo := repository.NewMockTodoRepository(t)
	emptyTodo := make([]entity.Todolist, 0)
	todoRepo.On("GetAll").Return(emptyTodo, nil)

	todoSvc := NewTodoService(todoRepo)
	endpoint := "/manage-todos"

	r := gin.New()
	r.GET(endpoint, todoSvc.TodolistHandlerGetAll)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the HTTP response status code
	assert.Equal(t, http.StatusOK, w.Code)

	respBody, err := io.ReadAll(w.Body)
	require.NoError(t, err)

	var result request.TodoResponseToGetAll
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("Failed to unmarshal JSON response body: %v", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, emptyTodo, result.Todos)
}

func testCreateSuccess(t *testing.T) {
	// Initialize mock repository
	todoRepo := repository.NewMockTodoRepository(t)

	// Set up mock behavior
	newTodo := &entity.Todolist{
		//ID:     1,
		Title:  "Makan",
		Status: false,
	}

	todoRepo.On("Create", "Makan").Return(newTodo, nil)

	// Initialize todo service with mock repository
	todoSvc := NewTodoService(todoRepo)

	// Call the create endpoint
	endpoint := "/manage-todo"
	r := gin.New()
	r.POST(endpoint, todoSvc.TodolistHandlerCreate)

	// Create an HTTP request to create a new Todo
	reqBody := bytes.NewBufferString(`{"title": "Makan"}`)
	req, err := http.NewRequest(http.MethodPost, endpoint, reqBody)
	require.NoError(t, err)

	// Send the request and read the response body
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	respBody, err := io.ReadAll(w.Body)
	if err != nil {
		log.Printf("Response body: %s\n", respBody)
	}
	//_ = fmt.Sprintf("%s", respBody)
	require.NoError(t, err)

	// Unmarshal the response body into a Todo object
	var result request.TodoResponse
	//str := fmt.Sprintf("%+v\n", responseBody)
	//fmt.Println(str)
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("Failed to unmarshal JSON response body: %v", err)
	}

	// Assert that the response has a 200 ok status code and returns the new Todo
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, newTodo, result.Data)
}

func testCreateInternalServerError(t *testing.T) {
	todoRepo := repository.NewMockTodoRepository(t)

	//newTodo := &entity.Todolist{
	//	Title:  "Makan",
	//	Status: false,
	//}
	//expectedErrors := errors.New("Internal Server Error")
	//todoRepo.On("Create", "Makan").Return(newTodo, expectedErrors)
	//
	//// Initialize todo service with mock repository
	//todoSvc := NewTodoService(todoRepo)
	//
	//// Call the create endpoint
	//endpoint := "/manage-todo"
	//r := gin.New()
	//r.POST(endpoint, todoSvc.TodolistHandlerCreate)
	//
	//// Create an HTTP request to create a new Todo
	//reqBody := bytes.NewBufferString(`{"title": "Makan"}`)
	//req, err := http.NewRequest(http.MethodPost, endpoint, reqBody)
	//require.NoError(t, err)
	//
	//w := httptest.NewRecorder()
	//r.ServeHTTP(w, req)
	//
	//respBody, err := io.ReadAll(w.Body)
	//require.NoError(t, err)
	//
	//var errResp respErr.ErrorResponse
	//err = json.Unmarshal(respBody, &errResp)
	//require.NoError(t, err)
	//
	//assert.Equal(t, http.StatusInternalServerError, w.Code)
	//assert.Equal(t, expectedErrors.Error(), errResp.Message)
	handler := NewTodoService(todoRepo)

	expectedError := errors.New("Internal Server Error")
	endpoint := "/manage-todo"

	todoRepo.On("Create", "Test Todo").Return(nil, expectedError)

	// Create valid input
	body := bytes.NewBufferString(`{"title": "Test Todo"}`)
	req, err := http.NewRequest("POST", endpoint, body)
	require.NoError(t, err)

	// Set up Gin context
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST(endpoint, handler.TodolistHandlerCreate)

	// Perform request
	c.Request = req
	r.ServeHTTP(w, req)

	respBody, err := io.ReadAll(w.Body)
	require.NoError(t, err)

	var errResp respErr.ErrorResponse
	err = json.Unmarshal(respBody, &errResp)
	require.NoError(t, err)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedError.Error(), errResp.Message)

	// Check mock call
	todoRepo.AssertCalled(t, "Create", "Test Todo")
}

func testCreateInvalid(t *testing.T) {
	// Setup
	todorepo := repository.NewMockTodoRepository(t)
	handler := NewTodoService(todorepo)

	expectedErrors := errors.New("Invalid input")

	//newTodo := &entity.Todolist{
	//	Title: "",
	//}
	//
	//todorepo.On("Create", newTodo.Title).Return(newTodo, expectedErrors).Once()

	endpoint := "/manage-todo"

	// Create invalid input
	body := bytes.NewBufferString(`{"title": ""}`)
	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	require.NoError(t, err)

	// Set up Gin context
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST(endpoint, handler.TodolistHandlerCreate)

	// Perform request
	c.Request = req
	r.ServeHTTP(w, req)

	respBody, err := io.ReadAll(w.Body)
	require.NoError(t, err)

	var errResp respErr.ErrorResponse
	err = json.Unmarshal(respBody, &errResp)
	require.NoError(t, err)

	// Check response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expectedErrors.Error(), errResp.Message)

	// Check mock call
	todorepo.AssertNotCalled(t, "Called", mock.Anything)
}

func testGetById(t *testing.T) {

}
