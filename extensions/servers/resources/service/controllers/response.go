package controllers

type ServiceResponse struct {
	Status bool `json:"status"`
	Count int `json:"count"`
	Data interface{} `json:"data"`
}
