package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"todoGin/model/entity"
	"todoGin/model/respErr"
	"todoGin/repository"
)

func TestTodolist(t *testing.T) {
	t.Run("testTodoServiceListFailedInternalServerError", testTodoServiceListFailedInternalServerError)
	t.Run("testGetAllSuccess", testGetAllSuccess)
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

	var errResp respErr.ErrorResponse
	err = json.Unmarshal(respBody, &errResp)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)

}

func testTodoServiceListFailedInternalServerError(t *testing.T) {
	todoRepo := repository.NewMockTodoRepository(t)
	expectedErr := errors.New("Internal Server Error")
	todoRepo.On("GetAll").Return(nil, expectedErr)

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
