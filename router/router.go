package router

import (
	"github.com/gin-gonic/gin"
	"todoGin/middleware"
	todoservice "todoGin/service"
)

type RouteBuilder struct {
	todoService *todoservice.Handler
}

func NewRouteBuilder(todoService *todoservice.Handler) *RouteBuilder {
	return &RouteBuilder{todoService: todoService}
}

func (rb *RouteBuilder) RouteInit() *gin.Engine {

	r := gin.New()
	//r.Use(gin.Logger())
	r.Use(gin.Recovery(), middleware.Logger(), middleware.XAPIKEY())

	r.GET("/manage-todos", rb.todoService.TodolistHandlerGetAll)
	r.POST("/manage-todo", rb.todoService.TodolistHandlerCreate)
	r.GET("/manage-todo/todo/:id", rb.todoService.TodolistHandlerGetByID)
	r.PUT("/manage-todo/todo/:id", rb.todoService.TodolistHandlerUpdate)
	r.DELETE("/manage-todo/todo/:id", rb.todoService.TodolistHandlerDelete)

	//err := r.Run()
	//if err != nil {
	//	return nil
	//}
	return r
}
