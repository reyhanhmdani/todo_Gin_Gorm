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
	t.Run("testTodolistHandlerGetByID", testTodolistHandlerGetByID)
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

func testTodolistHandlerGetByID(t *testing.T) {
	// inisiasi mocking
	mockTodoRepo := repository.NewMockTodoRepository(t)

	// inisiasi handler
	handler := NewTodoService(mockTodoRepo)

	// testing success
	mockTodoRepo.On("GetByID", int64(1)).Return(&entity.Todolist{ID: 1, Title: "Test Todo"}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/manage-todo/todo/1", nil)
	router := gin.Default()
	router.GET("/manage-todo/todo/:id", handler.TodolistHandlerGetByID)
	router.ServeHTTP(w, req)

	respBody, err := io.ReadAll(w.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)

	var response request.TodoResponse
	err = json.Unmarshal(respBody, &response)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, response.Message, "Success Get Id")

	// testing not found
	mockTodoRepo.On("GetByID", int64(2)).Return(nil, nil)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/manage-todo/todo/2", nil)
	router = gin.Default()
	router.GET("/manage-todo/todo/:id", handler.TodolistHandlerGetByID)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var res respErr.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, "Not Found", res.Message)
	assert.Equal(t, http.StatusNotFound, res.Status)

	// testing internal server error
	mockTodoRepo.On("GetByID", int64(3)).Return(nil, errors.New("Internal Server Error"))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/manage-todo/todo/3", nil)
	router = gin.Default()
	router.GET("/manage-todo/todo/:id", handler.TodolistHandlerGetByID)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var internal respErr.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &internal)
	require.NoError(t, err)
	assert.Equal(t, "Internal Server Error", internal.Message)
	assert.Equal(t, http.StatusInternalServerError, internal.Status)
}

func TestDeleteSuccess(t *testing.T) {
	mockTodoRepo := repository.NewMockTodoRepository(t)
	handler := NewTodoService(mockTodoRepo)

	// Testing Success
	mockTodoRepo.On("Delete", int64(1)).Return(int64(1), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/manage-todo/todo/1", nil)
	router := gin.Default()
	router.DELETE("/manage-todo/todo/:id", handler.TodolistHandlerDelete)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp request.TodoDeleteResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Succes Delete", resp.Message)

	// Testing Not Found
	//mockTodoRepo.On("Delete", int64(2)).Return("/manage-todo/todo/2")

}
