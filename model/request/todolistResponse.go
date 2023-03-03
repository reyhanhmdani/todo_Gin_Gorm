package request

import "todoGin/model/entity"

type TodoResponse struct {
	Status  interface{}     `json:"status"`
	Message interface{}     `json:"message"`
	Data    entity.Todolist `json:"data"`
}

type TodoResponseToGetAll struct {
	Status interface{}       `json:"status"`
	Data   int               `json:"data"`
	Todos  []entity.Todolist `json:"todos"`
}

type TodoIDResponse struct {
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type TodoDeleteResponse struct {
	Status  interface{} `json:"status"`
	Message interface{} `json:"message"`
}

type TodoUpdateResponse struct {
	Status  int         `json:"status"`
	Message interface{} `json:"data"`
	Todos   interface{} `json:"todos"`
}
