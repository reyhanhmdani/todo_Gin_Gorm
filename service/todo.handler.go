package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"todoGin/model/request"
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
	}
	logrus.Info(200, " Success Get All Data")
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    len(todos),
		"todos":   todos,
		"message": "Success Get All",
	})
	return
}
func (h *Handler) TodolistHandlerCreate(ctx *gin.Context) {
	todolist := new(request.TodolistCreateRequest)
	err := ctx.ShouldBindJSON(todolist)
	if err != nil {
		logrus.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}
	if _, errCreate := h.TodoRepository.Create(todolist.Title); errCreate != nil {
		logrus.Error(errCreate)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed Create Data",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	//if errCreate := db.Raw("INSERT INTO todolist(title) values (?)").Scan(&todolist)
	logrus.Info(http.StatusOK, " Success Create Todo", todolist)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    todolist,
	})
	return
}
func (h *Handler) TodolistHandlerGetByID(ctx *gin.Context) {
	userId := ctx.Param("id")
	todoID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	todo, err := h.TodoRepository.GetByID(todoID)
	if err != nil {
		logrus.Errorf("failed when get todo by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	if todo == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "USER NOT FOUND",
			"status":  http.StatusNotFound,
		})
		return
	}
	logrus.Info(200, " Success Get By ID")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get ID",
		"data":    todo,
		"Status":  http.StatusOK,
	})
	return
}

func (h *Handler) TodolistHandlerUpdate(ctx *gin.Context) {
	userId := ctx.Param("id")
	todoID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	reqBody := new(request.TodolistUpdateRequest)
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}
	ErrId, err := h.TodoRepository.GetByID(todoID)
	if err != nil {
		logrus.Errorf("failed when get todo by id: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	if ErrId == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "USER NOT FOUND",
			"status":  http.StatusNotFound,
		})
		return
	}
	rowsAffected, err := h.TodoRepository.Update(todoID, reqBody.ReqTodo())
	if err != nil {
		logrus.Errorf("failed when updating todo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if rowsAffected == nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
		})
		return
	}

	logrus.Info(200, " Success Update Todo")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully update Data",
		"status":  http.StatusOK,
		"data":    reqBody,
	})
	return

}
func (h *Handler) TodolistHandlerDelete(ctx *gin.Context) {
	userId := ctx.Param("id")
	todoID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  400,
		})
		return
	}
	isFound, err := h.TodoRepository.Delete(todoID)
	if err != nil {
		logrus.Errorf("failed when deleting todo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	//fmt.Println(isFound)
	if isFound == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "NOT FOUND ID",
			"status":  404,
		})
		return
	}
	logrus.Info(200, " Success DELETE")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully Delete Data",
		"status":  http.StatusOK,
	})

	//fmt.Println()
	return
}
