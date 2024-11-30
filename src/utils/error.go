package utils

import (
	"cij_api/src/model"
	"fmt"
)

type Error struct {
	Message string        `json:"message"`
	Code    string        `json:"code"`
	Fields  []model.Field `json:"fields,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) GetCode() string {
	return e.Code
}

func (e *Error) SetCode(code string) {
	e.Code = code
}

func (e Error) GetFields() []model.Field {
	return e.Fields
}

func NewErrorCode(errorType ErrorType, errorEntity ErrorEntity, identifier string) string {
	return fmt.Sprintf("%d%d%s", errorType, errorEntity, identifier)
}

func NewError(message string, code string) Error {
	return Error{
		Message: message,
		Code:    code,
	}
}

func NewErrorWithFields(message string, code string, fields []model.Field) Error {
	return Error{
		Message: message,
		Code:    code,
		Fields:  fields,
	}
}

// Error code
type ErrorType int

const (
	ValidationErrorCode ErrorType = 1
	DatabaseErrorCode   ErrorType = 2
	ServiceErrorCode    ErrorType = 3
	ControllerErrorCode ErrorType = 4
)

type ErrorEntity int

const (
	UserErrorType       ErrorEntity = 1
	PersonErrorType     ErrorEntity = 2
	AddressErrorType    ErrorEntity = 3
	DisabilityErrorType ErrorEntity = 4
	CompanyErrorType    ErrorEntity = 5
	NewsErrorType       ErrorEntity = 6
	ConfigErrorType     ErrorEntity = 7
	ActivityErrorType   ErrorEntity = 8
	ReportsErrorType    ErrorEntity = 9
	VacancyErrorType    ErrorEntity = 10
)
