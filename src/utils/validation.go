package utils

import "cij_api/src/model"

func ValidateAddress(addressRequest model.AddressRequest) Error {
	fieldsWithErrors := []model.Field{}

	if addressRequest.Street == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "street"})
	}

	if addressRequest.Number == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "number"})
	}

	if addressRequest.Neighborhood == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "neighborhood"})
	}

	if addressRequest.City == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "city"})
	}

	if addressRequest.State == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "state"})
	}

	if len(addressRequest.State) != 2 {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "state", Value: "state must have 2 characters"})
	}

	if addressRequest.ZipCode == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "zip code"})
	}

	if len(addressRequest.ZipCode) != 8 {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "zip code", Value: "zip code must have 8 digits"})
	}

	if addressRequest.Country == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "country"})
	}

	if len(fieldsWithErrors) > 0 {
		errorCode := NewErrorCode(ValidationErrorCode, AddressErrorType, "01")

		return NewErrorWithFields("required fields are missing", errorCode, fieldsWithErrors)
	}

	return Error{}
}

func ValidateUser(user model.UserRequest) Error {
	fieldsWithErrors := []model.Field{}

	if user.Email == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "email"})
	}

	if user.Password == "" {
		fieldsWithErrors = append(fieldsWithErrors, model.Field{Name: "password"})
	}

	if len(fieldsWithErrors) > 0 {
		errorCode := NewErrorCode(ValidationErrorCode, UserErrorType, "01")

		return NewErrorWithFields("required fields are missing", errorCode, fieldsWithErrors)
	}

	return Error{}
}
