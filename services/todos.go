package services

import (
	"github.com/melardev/GoMuxGormApiCrudPagination/infrastructure"
	"github.com/melardev/GoMuxGormApiCrudPagination/models"
)

func FetchTodos(page, pageSize int) ([]models.Todo, int) {
	var todos []models.Todo
	totalTodosCount := 0
	database := infrastructure.GetDb()
	database.Model(&models.Todo{}).Count(&totalTodosCount)
	database.Select("id, title, completed, created_at, updated_at").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at desc").
		Find(&todos)

	return todos, totalTodosCount
}

func FetchPendingTodos(page, pageSize int, completed bool) ([]models.Todo, int) {
	var todos []models.Todo
	var totalTodosCount int
	database := infrastructure.GetDb()
	database.Model(&models.Todo{}).
		// Warning: This line would not work: http://gorm.io/docs/query.html#Struct-amp-Map
		// Where(&models.Todo{Completed: completed}).
		// Use this instead
		Where("completed = ?", completed).
		Count(&totalTodosCount)
	database.Select("id, title, completed, created_at, updated_at").
		Where("completed = ?", completed).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at desc").
		Find(&todos)

	return todos, totalTodosCount
}

func DeleteAllTodos() {
	database := infrastructure.GetDb()
	database.Model(&models.Todo{}).Delete(&models.Todo{})
}

func FetchById(id uint) (todo models.Todo, err error) {
	database := infrastructure.GetDb()
	err = database.Model(&models.Todo{}).First(&todo, id).Error
	return
}

func DeleteTodo(todo *models.Todo) error {
	database := infrastructure.GetDb()
	return database.Delete(todo).Error
}

func CreateTodo(title, description string, completed bool) (todo models.Todo, err error) {
	database := infrastructure.GetDb()
	todo = models.Todo{Title: title, Description: description, Completed: completed}
	err = database.Create(&todo).Error
	return todo, err
}

func UpdateTodo(id uint, title, description string, completed bool) (todo models.Todo, err error) {
	todo, err = FetchById(id)
	if err != nil {
		return
	}
	todo.Title = title

	// TODO: handle this in a better way, the user should be able to set description to empty string
	// The intention is to check against nil but in go there are no nil strings, so we can not know
	// if the user intended to udpate the description to empty string or just update the other fields other than description.
	if description != "" {
		todo.Description = description
	}

	todo.Completed = completed
	database := infrastructure.GetDb()

	// database.Model(&todo).Updates(map[string]interface{}{"title": title, "description": description, "completed": completed})
	// database.Model(&todo).Updates(models.Todo{Title: title, Description: description, Completed: completed})

	database.Save(&todo)

	return
}
