package dtos

import (
	"fmt"
	"github.com/melardev/GoMuxGormApiCrudPagination/models"
	"gopkg.in/go-playground/validator.v8"
	"math"
	"net/http"
)

func CreatePagedResponse(request *http.Request, resources []interface{}, resourceName string, page, pageSize, totalItemsCount int) map[string]interface{} {
	response := CreatePageMeta(request, len(resources), page, pageSize, totalItemsCount)
	response[resourceName] = resources
	return response
}

func CreatePageMeta(request *http.Request, loadedItemsCount, page, page_size, totalItemsCount int) map[string]interface{} {
	page_meta := map[string]interface{}{}
	page_meta["offset"] = (page - 1) * page_size
	page_meta["requested_page_size"] = page_size
	page_meta["current_page_number"] = page
	page_meta["current_items_count"] = loadedItemsCount
	page_meta["total_items_count"] = totalItemsCount

	page_meta["prev_page_number"] = 1
	numberOfPages := int(math.Ceil(float64(totalItemsCount) / float64(page_size)))
	page_meta["number_of_pages"] = numberOfPages

	if page < numberOfPages {
		page_meta["has_next_page"] = true
		page_meta["next_page_number"] = page + 1
	} else {
		page_meta["has_next_page"] = false
		page_meta["next_page_number"] = 1
	}
	if page > 1 {
		page_meta["prev_page_number"] = page - 1
		page_meta["has_prev_page"] = true
	} else {
		page_meta["has_prev_page"] = false
		page_meta["prev_page_number"] = 1
	}

	page_meta["next_page_url"] = fmt.Sprintf("%v?page=%d&page_size=%d", request.URL.Path, page_meta["next_page_number"], page_meta["requested_page_size"])
	page_meta["prev_page_url"] = fmt.Sprintf("%s?page=%d&page_size=%d", request.URL.Path, page_meta["prev_page_number"], page_meta["requested_page_size"])

	response := map[string]interface{}{
		"success":   true,
		"page_meta": page_meta,
	}

	return response
}

func CreateBadRequestErrorDto(err error) map[string]interface{} {
	res := make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	res["full_messages"] = make([]string, len(errs))
	count := 0
	for _, v := range errs {
		if v.ActualTag == "required" {
			var message = fmt.Sprintf("%v is required", v.Field)
			res["full_messages"].([]string)[count] = message
		} else {
			var message = fmt.Sprintf("%v has to be %v", v.Field, v.ActualTag)
			res["full_messages"].([]string)[count] = message
		}
		count++
	}
	return res
}

func CreateErrorDtoWithMessage(message string) map[string]interface{} {
	return map[string]interface{}{
		"success":       false,
		"full_messages": []string{message},
	}
}

func GetTodoDetaislDto(todo *models.Todo) map[string]interface{} {
	response := GetTodoDto(todo, true)
	return response
}

func CreateSuccessWithMessageDto(message string) interface{} {
	return CreateSuccessWithMessagesDto([]string{message})
}

func CreateSuccessWithMessagesDto(messages []string) interface{} {
	return map[string]interface{}{
		"success":       true,
		"full_messages": messages,
	}
}

func CreateSuccessWithDtoAndMessagesDto(data map[string]interface{}, messages []string) map[string]interface{} {
	data["success"] = true
	data["full_messages"] = messages
	return data
}

func CreateSuccessWithDtoAndMessageDto(data map[string]interface{}, message string) map[string]interface{} {
	return CreateSuccessWithDtoAndMessagesDto(data, []string{message})
}
