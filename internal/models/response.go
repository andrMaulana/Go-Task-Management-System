package models

type Meta struct {
	Code    int    `json:"code"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}
