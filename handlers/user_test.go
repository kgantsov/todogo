package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/kgantsov/todogo/models"
	"net/http"
	"testing"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{ID: 1, Name: "Mike", Email: "mike@gmail.com", Password: "111"},
	{ID: 2, Name: "Ben", Email: "ben@gmail.com", Password: "111"},
	{ID: 3, Name: "Kevin", Email: "kevin@gmail.com", Password: "111"},
	{ID: 4, Name: "Tom", Email: "tom@gmail.com", Password: "111"},
	{ID: 5, Name: "Oliver", Email: "oliver@gmail.com", Password: "111"},
	{ID: 6, Name: "Pol", Email: "pol@gmail.com", Password: "111"},
}

func CreateUserFixtures(db *gorm.DB) {
	for _, user := range users {
		db.Create(&user)
	}
}

func TestGetUsers(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/user/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res []models.User

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range users {
		if res[i].ID != users[i].ID {
			t.Errorf("Response body should be `%s`, was  %s", users[i].ID, res[i].ID)
		}
		if res[i].Name != users[i].Name {
			t.Errorf("Response body should be `%s`, was  %s", users[i].Name, res[i].Name)
		}
		if res[i].Email != users[i].Email {
			t.Errorf("Response body should be `%s`, was  %s", users[i].Email, res[i].Email)
		}
		if res[i].Password != users[i].Password {
			t.Errorf("Response body should be `%s`, was  %s", users[i].Password, res[i].Password)
		}
	}
}

func TestGetUsersEmptyTable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/user/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestGetUser(t *testing.T) {
	user := users[1]

	req, _ := http.NewRequest("GET", "/api/v1/user/2/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res models.User

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.ID != user.ID {
		t.Errorf("Response body should be `%s`, was  %s", user.ID, res.ID)
	}
	if res.Name != user.Name {
		t.Errorf("Response body should be `%s`, was  %s", user.Name, res.Name)
	}
	if res.Email != user.Email {
		t.Errorf("Response body should be `%s`, was  %s", user.Email, res.Email)
	}
	if res.Password != user.Password {
		t.Errorf("Response body should be `%s`, was  %s", user.Password, res.Password)
	}
}

func TestGetUserWrongID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/user/777/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusNotFound, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["error"] != "User not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'User not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestCreateUser(t *testing.T) {
	var jsonStr = []byte(`{"name": "Petr", "email": "petr@gmail.com", "password": "222"}`)

	req, _ := http.NewRequest("POST", "/api/v1/user/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusCreated, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res models.User

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.Name != "Petr" {
		t.Errorf("Response body should be `Petr`, was  %s", res.Name)
	}
	if res.Email != "petr@gmail.com" {
		t.Errorf("Response body should be `petr@gmail.com`, was  %s", res.Email)
	}
	if res.Password != hashPassword("222") {
		t.Errorf("Response body should be `%s`, was  %s", hashPassword("222"), res.Password)
	}
}

func TestCreateUserMissedFields(t *testing.T) {
	var jsonStr = []byte(`{}`)

	req, _ := http.NewRequest("POST", "/api/v1/user/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusBadRequest, resp.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	var jsonStr = []byte(`{"name": "Tim", "email": "tim@gmail.com", "password": "444"}`)

	req, _ := http.NewRequest("PUT", "/api/v1/user/2/", bytes.NewBuffer(jsonStr))
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res models.User

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.ID != 2 {
		t.Errorf("Response body should be `%s`, was  %s", 2, res.ID)
	}
	if res.Name != "Tim" {
		t.Errorf("Response body should be `%s`, was  %s", "Tim", res.Name)
	}
	if res.Email != "tim@gmail.com" {
		t.Errorf("Response body should be `%s`, was  %s", "tim@gmail.com", res.Email)
	}
	if res.Password != hashPassword("444") {
		t.Errorf("Response body should be `%s`, was  %s", hashPassword("444"), res.Password)
	}
}

func TestUpdateUserWrongID(t *testing.T) {
	var jsonStr = []byte(`{"name": "Tim", "email": "tim@gmail.com", "password": "444"}`)

	req, _ := http.NewRequest("PUT", "/api/v1/user/777/", bytes.NewBuffer(jsonStr))
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusNotFound, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["error"] != "User not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'User not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestDeleteUser(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/user/2/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusNoContent, resp.Code)
	}
}

func TestDeleteUserWrongID(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/user/777/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusNotFound, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["error"] != "User not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'User not found'. Got '%s'",
			res["error"],
		)
	}
}
