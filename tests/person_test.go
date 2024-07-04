package tests

import (
	"bytes"
	"cij_api/src/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

const BASE_URL = "http://localhost:3040"

var personId int = 1
var token string

func TestCreatePerson(t *testing.T) {
	url := BASE_URL + "/people"

	personJson, err := os.ReadFile("./mocks/person_request.json")
	if err != nil {
		t.Error(err)
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(personJson))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error("Expected status code 200, got", response.StatusCode)
	}

	defer response.Body.Close()
}

func TestListPeople(t *testing.T) {
	url := BASE_URL + "/people"

	response, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error("Expected status code 200, got", response.StatusCode)
	}

	body, _ := io.ReadAll(response.Body)

	var responseData model.ResponseData[[]model.PersonResponse]
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		t.Error(err)
	}

	if len(responseData.Data) == 0 {
		t.Error("Expected at least one person, got", len(responseData.Data))
	}

	personId = responseData.Data[len(responseData.Data)-1].Id

	defer response.Body.Close()
}

func TestGetPerson(t *testing.T) {
	url := BASE_URL + "/people/" + fmt.Sprint(personId)

	response, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		t.Error("Expected status code 200 or 404, got", response.StatusCode)
	}

	defer response.Body.Close()
}

func TestLoginPerson(t *testing.T) {
	url := BASE_URL + "/login"

	personJson, err := os.ReadFile("./mocks/person_login.json")
	if err != nil {
		t.Error(err)
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(personJson))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error("Expected status code 200, got", response.StatusCode)
	}

	body, _ := io.ReadAll(response.Body)
	var responseData model.LoginResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		t.Error(err)
	}

	token = responseData.Token
	if token == "" {
		t.Error("Expected token, got", token)
	}

	defer response.Body.Close()
}

func TestUpdatePerson(t *testing.T) {
	url := BASE_URL + "/people/" + fmt.Sprint(personId)

	personJson, err := os.ReadFile("./mocks/person_update.json")
	if err != nil {
		t.Error(err)
	}

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(personJson))
	if err != nil {
		t.Error(err)
	}

	request.Header = http.Header{
		"Authorization": {token},
		"Content-Type":  {"application/json"},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error("Expected status code 200, got", response.StatusCode)
	}

	defer response.Body.Close()
}

func TestDeletePerson(t *testing.T) {
	url := BASE_URL + "/people/" + fmt.Sprint(personId)

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Error(err)
	}

	request.Header = http.Header{
		"Authorization": {token},
		"Content-Type":  {"application/json"},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error("Expected status code 200, got", response.StatusCode)
	}

	defer response.Body.Close()
}
