package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kgantsov/todogo/models"
)

var testDBConnectionString = "postgresql://root@localhost:26257/todogo_test?sslmode=disable"

func TestLogin(t *testing.T) {
	var jsonStr = []byte(`{"email": "mike@gmail.com", "password": "111"}`)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDbURI(testDBConnectionString, false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusCreated, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	req, _ = http.NewRequest("GET", "/api/v1/list/", nil)

	req.Header.Set("Auth-Token", res["token"])

	resp = ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString = resp.Body.String()

	var todoList []models.TodoList

	err = json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(todoList) != 3 {
		t.Errorf("Response body should be `3`, was  %d", len(todoList))
	}

	if todoList[0].ID != todoLists[0].ID {
		t.Errorf("Response body should be `1`, was  %s", todoList[0].ID)
	}
	if todoList[0].Title != "Shopping list" {
		t.Errorf("Response body should be `Shopping list`, was  %s", todoList[0].Title)
	}

	if todoList[1].ID != todoLists[0].ID {
		t.Errorf("Response body should be `2`, was  %s", todoList[1].ID)
	}
	if todoList[1].Title != "Work list" {
		t.Errorf("Response body should be `Work list`, was  %s", todoList[1].Title)
	}

	if todoList[2].ID != todoLists[0].ID {
		t.Errorf("Response body should be `3`, was  %s", todoList[2].ID)
	}
	if todoList[2].Title != "Sport list" {
		t.Errorf("Response body should be `Sport list`, was  %s", todoList[2].Title)
	}
}

func TestLoginNonExistentUser(t *testing.T) {
	var jsonStr = []byte(`{"email": "mike@gmail.com", "password": "111"}`)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login/", bytes.NewBuffer(jsonStr))

	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusUnauthorized, resp.Code)
	}
}

func TestLoginNewUser(t *testing.T) {
	var jsonStr = []byte(`{"name": "Ilon", "email": "ilon@gmail.com", "password": "555"}`)

	req, _ := http.NewRequest("POST", "/api/v1/user/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDbURI(testDBConnectionString, false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusCreated, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var user models.User

	err := json.Unmarshal([]byte(bodyAsString), &user)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Name != "Ilon" {
		t.Errorf("Response body should be `Ilon`, was  %s", user.Name)
	}
	if user.Email != "ilon@gmail.com" {
		t.Errorf("Response body should be `ilon@gmail.com`, was  %s", user.Email)
	}
	if user.Password != hashPassword("555") {
		t.Errorf("Response body should be `%s`, was  %s", hashPassword("555"), user.Password)
	}

	jsonStr = []byte(`{"email": "ilon@gmail.com", "password": "555"}`)

	req, _ = http.NewRequest("POST", "/api/v1/auth/login/", bytes.NewBuffer(jsonStr))

	resp = ExecuteRequest(db, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusCreated, resp.Code)
	}

	bodyAsString = resp.Body.String()

	var res map[string]string

	err = json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	req, _ = http.NewRequest("GET", "/api/v1/list/", nil)

	req.Header.Set("Auth-Token", res["token"])

	resp = ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString = resp.Body.String()

	err = json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func TestLoginWithoutToken(t *testing.T) {
	var jsonStr = []byte(`{"email": "mike@gmail.com", "password": "111"}`)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login/", bytes.NewBuffer(jsonStr))

	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusCreated, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	req, _ = http.NewRequest("GET", "/api/v1/list/", nil)

	resp = ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestLoginIncorrectPassword(t *testing.T) {
	var jsonStr = []byte(`{"email": "mike@gmail.com", "password": "222"}`)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login/", bytes.NewBuffer(jsonStr))

	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusUnauthorized, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["error"] != "Login or password is incorrect" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Login or password is incorrect'. Got '%s'",
			res["error"],
		)
	}
}
