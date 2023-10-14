package utils

import "go.mongodb.org/mongo-driver/bson"

type apiResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseSuccess(data interface{}, message string) apiResponse {
	resp := apiResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
	return resp
}

func ResponseError(data interface{}, message string) apiResponse {
	resp := apiResponse{
		Status:  false,
		Message: message,
		Data:    data,
	}
	return resp
}

type paginationCustom struct {
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalCount int         `json:"total_count"`
	Data       interface{} `json:"data"`
}

func PaginationCustom(page int, perPage int, total int, data interface{}) paginationCustom {
	resp := paginationCustom{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		Data:       data,
	}

	return resp
}

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
