package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/melardev/GoMuxGormApiCrudPagination/dtos"
	"github.com/melardev/GoMuxGormApiCrudPagination/models"
	"github.com/melardev/GoMuxGormApiCrudPagination/services"
	"net/http"
	"strconv"
)

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	page, pageSize := getPagingParams(r)
	todos, totalTodoCount := services.FetchTodos(page, pageSize)
	sendAsJson(w, http.StatusOK, dtos.CreateTodoPagedResponse(r, todos, page, pageSize, totalTodoCount))
}

func GetAllPendingTodos(w http.ResponseWriter, r *http.Request) {
	page, pageSize := getPagingParams(r)
	todos, totalTodoCount := services.FetchPendingTodos(page, pageSize, false)
	sendAsJson(w, http.StatusOK, dtos.CreateTodoPagedResponse(r, todos, page, pageSize, totalTodoCount))
}
func GetAllCompletedTodos(w http.ResponseWriter, r *http.Request) {
	page, pageSize := getPagingParams(r)
	todos, totalTodoCount := services.FetchPendingTodos(page, pageSize, true)
	sendAsJson(w, http.StatusOK, dtos.CreateTodoPagedResponse(r, todos, page, pageSize, totalTodoCount))
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	id64, _ := strconv.ParseUint(id, 10, 32)
	todo, err := services.FetchById(uint(id64))
	if err != nil {
		sendAsJson(w, http.StatusNotFound, dtos.CreateErrorDtoWithMessage("Could not find Todo"))
		return
	}

	// Just to prove that sendAsJson2 also works
	sendAsJson2(w, http.StatusOK, dtos.GetTodoDetaislDto(&todo))
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	todo := models.Todo{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&todo); err != nil {
		sendAsJson(w, http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}
	defer r.Body.Close()

	todo, err := services.CreateTodo(todo.Title, todo.Description, todo.Completed)
	if err != nil {
		sendAsJson(w, http.StatusInternalServerError, dtos.CreateErrorDtoWithMessage(err.Error()))
		return
	}

	sendAsJson(w, http.StatusCreated, dtos.CreateTodoCreatedDto(&todo))
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendAsJson(w, http.StatusBadRequest, dtos.CreateErrorDtoWithMessage("You must set an ID"))
		return
	}

	var todoInput models.Todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&todoInput); err != nil {
		sendAsJson(w, http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}
	defer r.Body.Close()

	todo, err := services.UpdateTodo(uint(id), todoInput.Title, todoInput.Description, todoInput.Completed)
	if err != nil {
		sendAsJson(w, http.StatusInternalServerError, dtos.CreateErrorDtoWithMessage(err.Error()))
		return
	}

	sendAsJson(w, http.StatusOK, dtos.CreateTodoUpdatedDto(&todo))
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendAsJson(w, http.StatusBadRequest, dtos.CreateErrorDtoWithMessage("You must set an ID"))
		return
	}
	todo, err := services.FetchById(uint(id))
	if err != nil {
		sendAsJson(w, http.StatusNotFound, dtos.CreateErrorDtoWithMessage("todo not found"))
		return
	}

	err = services.DeleteTodo(&todo)

	if err != nil {
		sendAsJson(w, http.StatusNotFound, dtos.CreateErrorDtoWithMessage("Could not delete Todo"))
		return
	}

	sendAsJson(w, http.StatusNoContent, dtos.CreateSuccessWithMessageDto("Todo deleted successfully"))
}

func DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	services.DeleteAllTodos()
	sendAsJson(w, http.StatusNoContent, dtos.CreateSuccessWithMessageDto("All Todos deleted successfully"))
}
