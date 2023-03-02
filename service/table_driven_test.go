package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"todoGin/model/request"
	"todoGin/model/respErr"
	"todoGin/repository"
)

func TestTodolistHandlerDelete(t *testing.T) {
	mockRepo := repository.NewMockTodoRepository(t)
	handler := NewTodoService(mockRepo)
	gin.SetMode(gin.TestMode)

	// Test cases
	testCases := []struct {
		name      string
		todoID    int64
		isFound   int64
		expStatus int
		expResp   interface{}
	}{
		{
			name:      "Success",
			todoID:    1,
			isFound:   1,
			expStatus: http.StatusOK,
			expResp: request.TodoDeleteResponse{
				Status:  http.StatusOK,
				Message: "Success Delete",
			},
		},
		{
			name:      "Not Found",
			todoID:    2,
			isFound:   0,
			expStatus: http.StatusNotFound,
			expResp: respErr.ErrorResponse{
				Message: "Not Found",
				Status:  http.StatusNotFound,
			},
		},
		//{
		//	name:      "Internal Server Error",
		//	todoID:    3,
		//	isFound:   0,
		//	expStatus: http.StatusInternalServerError,
		//	expResp: respErr.ErrorResponse{
		//		Message: "Internal Server Error",
		//		Status:  http.StatusInternalServerError,
		//	},
		//},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("Delete", tc.todoID).Return(tc.isFound, nil)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodDelete, "/manage-todo/todo/"+strconv.FormatInt(tc.todoID, 10), nil)
			router := gin.Default()
			router.DELETE("/manage-todo/todo/:id", handler.TodolistHandlerDelete)
			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expStatus, w.Code)

			var respBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &respBody)
			if err != nil {
				log.Print(err)
			}
			require.NoError(t, err)
			assert.IsEqual(t, reflect.DeepEqual(tc.expResp, &respBody))

			mockRepo.AssertExpectations(t)
		})
	}
}
