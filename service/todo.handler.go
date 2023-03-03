package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"todoGin/model/request"
	"todoGin/model/respErr"
	"todoGin/repository"
)

type Handler struct {
	TodoRepository repository.TodoRepository
}

func NewTodoService(todoRepo repository.TodoRepository) *Handler {
	return &Handler{
		TodoRepository: todoRepo,
	}
}

func (h *Handler) TodolistHandlerGetAll(ctx *gin.Context) {
	todos, err := h.TodoRepository.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &respErr.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}
	logrus.Info(http.StatusOK, " Success Get All Data")
	//ctx.AbortWithStatusJSON(http.StatusOK, todos)
	ctx.AbortWithStatusJSON(http.StatusOK, request.TodoResponseToGetAll{
		Status: "Success Get All",
		Data:   len(todos),
		Todos:  todos,
	})
	return
}
func (h *Handler) TodolistHandlerCreate(ctx *gin.Context) {
	todolist := new(request.TodolistCreateRequest)
	err := ctx.ShouldBindJSON(todolist)
	if err != nil {
		logrus.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.ErrorResponse{
			Message: "Invalid input",
			Status:  http.StatusBadRequest,
		})
		return
	}
	newTodo, errCreate := h.TodoRepository.Create(todolist.Title)
	if errCreate != nil {
		logrus.Error(errCreate)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}

	logrus.Info(http.StatusOK, " Success Create Todo", todolist)
	ctx.JSON(http.StatusOK, request.TodoResponse{
		Status:  http.StatusOK,
		Message: "New Todo Created",
		Data:    *newTodo,
	})
	return
}
func (h *Handler) TodolistHandlerGetByID(ctx *gin.Context) {
	userId := ctx.Param("id")
	todoID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.ErrorResponse{
			Message: "Bad request",
			Status:  http.StatusBadRequest,
		})
		return
	}
	todo, err := h.TodoRepository.GetByID(todoID)
	if err != nil {
		logrus.Errorf("failed when get todo by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}
	if todo == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, respErr.ErrorResponse{
			Message: "Not Found",
			Status:  http.StatusNotFound,
		})
		return
	}
	logrus.Info(http.StatusOK, " Success Get By ID")
	ctx.JSON(http.StatusOK, request.TodoResponse{
		Status:  http.StatusOK,
		Message: "Success Get Id",
		Data:    *todo,
	})
	return
}

func (h *Handler) TodolistHandlerUpdate(ctx *gin.Context) {
	userId := ctx.Param("id")
	todoID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.ErrorResponse{
			Message: "parse ID error",
			Status:  http.StatusBadRequest,
		})
		return
	}
	reqBody := new(request.TodolistUpdateRequest)
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.ErrorResponse{
			Message: "Bad request",
			Status:  http.StatusBadRequest,
		})
		return
	}
	ErrId, err := h.TodoRepository.GetByID(todoID)
	if err != nil {
		logrus.Errorf("failed when get todo by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}
	if ErrId == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, respErr.ErrorResponse{
			Message: "ID not Found",
			Status:  http.StatusNotFound,
		})
		return
	}
	rowsAffected, err := h.TodoRepository.Update(todoID, reqBody.ReqTodo())
	if err != nil {
		logrus.Errorf("failed when updating todo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}
	if rowsAffected == nil {
		ctx.AbortWithStatusJSON(http.StatusOK, request.TodoIDResponse{
			Message: "Not Change",
			Data:    reqBody,
		})
		return
	}

	logrus.Info(http.StatusOK, " Success Update Todo")
	ctx.JSON(http.StatusOK, request.TodoUpdateResponse{
		Status:  http.StatusOK,
		Message: "Success Update Todo",
		Todos:   reqBody,
	})
	return

}
func (h *Handler) TodolistHandlerDelete(ctx *gin.Context) {
	userId := ctx.Param("id")
	todoID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respErr.ErrorResponse{
			Message: "Parse ID Error",
			Status:  http.StatusBadRequest,
		})
		return
	}
	isFound, err := h.TodoRepository.Delete(todoID)
	if err != nil {
		logrus.Errorf("failed when deleting todo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
		return
	}
	//fmt.Println(isFound)
	if isFound == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, respErr.ErrorResponse{
			Message: "Not Found",
			Status:  http.StatusNotFound,
		})
		return
	}
	logrus.Info(http.StatusOK, " Success DELETE")
	ctx.JSON(http.StatusOK, request.TodoDeleteResponse{
		Status:  http.StatusOK,
		Message: "Success Delete",
	})
	return
}
