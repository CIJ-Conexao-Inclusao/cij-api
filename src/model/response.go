package model

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type Response struct {
	Message string  `json:"message"`
	Code    string  `json:"code,omitempty"`
	Fields  []Field `json:"fields,omitempty"`
}

type ResponseData[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type LoginResponse struct {
	Token    string      `json:"token,omitempty"`
	Code     string      `json:"code,omitempty"`
	UserInfo interface{} `json:"user_info,omitempty"`
	Message  string      `json:"message,omitempty"`
}
