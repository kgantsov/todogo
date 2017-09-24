package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/models"
)

var users = []models.User{
	{ID: 1, Name: "Mike", Email: "mike@gmail.com", Password: hashPassword("111")},
	{ID: 2, Name: "Ben", Email: "ben@gmail.com", Password: hashPassword("111")},
	{ID: 3, Name: "Kevin", Email: "kevin@gmail.com", Password: hashPassword("111")},
	{ID: 4, Name: "Tom", Email: "tom@gmail.com", Password: hashPassword("111")},
	{ID: 5, Name: "Oliver", Email: "oliver@gmail.com", Password: hashPassword("111")},
	{ID: 6, Name: "Pol", Email: "pol@gmail.com", Password: hashPassword("111")},
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

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

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

	if res[0].ID != users[0].ID {
		t.Errorf("Response body should be `%s`, was  %s", users[0].ID, res[0].ID)
	}
	if res[0].Name != users[0].Name {
		t.Errorf("Response body should be `%s`, was  %s", users[0].Name, res[0].Name)
	}
	if res[0].Email != users[0].Email {
		t.Errorf("Response body should be `%s`, was  %s", users[0].Email, res[0].Email)
	}
	if res[0].Password != users[0].Password {
		t.Errorf("Response body should be `%s`, was  %s", users[0].Password, res[0].Password)
	}
}

func TestGetUsersEmptyTable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/user/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestGetUsersEmptyTableNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/user/", nil)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestGetUser(t *testing.T) {
	user := users[2]

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/user/%d/", user.ID), nil)
	token, _ := createToken(user.ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

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
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/user/%d/", users[2].ID), nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["error"] != "Access denied" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Access denied'. Got '%s'",
			res["error"],
		)
	}
}

func TestGetUserNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/user/777/", nil)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestCreateUser(t *testing.T) {
	var jsonStr = []byte(`{"name": "Petr", "email": "petr@gmail.com", "password": "222"}`)

	req, _ := http.NewRequest("POST", "/api/v1/user/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

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

func TestCreateUserExisingEmail(t *testing.T) {
	var jsonStr = []byte(`{"name": "Mike", "email": "mike@gmail.com", "password": "222"}`)

	req, _ := http.NewRequest("POST", "/api/v1/user/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusConflict {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusConflict, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["error"] != "User with this email already exists" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'User with this email already exists'. Got '%s'",
			res["error"],
		)
	}
}

func TestCreateUserMissedFields(t *testing.T) {
	var jsonStr = []byte(`{}`)

	req, _ := http.NewRequest("POST", "/api/v1/user/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusBadRequest, resp.Code)
	}
}

func TestCreateUserInvalidEmail(t *testing.T) {
	var jsonStr = []byte(`{"name": "Petr", "email": "invalidemail.com", "password": "222"}`)

	req, _ := http.NewRequest("POST", "/api/v1/user/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusBadRequest, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["error"] != "Email address is not valid" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Access denied'. Got '%s'",
			res["error"],
		)
	}
}

func TestUpdateUser(t *testing.T) {
	user := users[2]
	var jsonStr = []byte(`{"name": "Tim", "email": "tim@gmail.com", "password": "444"}`)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/user/%d/", user.ID), bytes.NewBuffer(jsonStr))
	token, _ := createToken(user.ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

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

	req, _ := http.NewRequest(
		"PUT", fmt.Sprintf("/api/v1/user/%d/", users[2].ID), bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["error"] != "Access denied" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Access denied'. Got '%s'",
			res["error"],
		)
	}
}

func TestUpdateUserNoAuthToken(t *testing.T) {
	var jsonStr = []byte(`{"name": "Tim", "email": "tim@gmail.com", "password": "444"}`)

	req, _ := http.NewRequest("PUT", "/api/v1/user/777/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestUpdateUserInvalidEmail(t *testing.T) {
	user := users[2]
	var jsonStr = []byte(`{"name": "Tim", "email": "invalidemail.com", "password": "444"}`)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/user/%d/", user.ID), bytes.NewBuffer(jsonStr))
	token, _ := createToken(user.ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["error"] != "Email address is not valid" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Access denied'. Got '%s'",
			res["error"],
		)
	}
}

func TestDeleteUser(t *testing.T) {
	user := users[2]
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/user/%d/", user.ID), nil)
	token, _ := createToken(user.ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusNoContent, resp.Code)
	}
}

func TestDeleteUserWrongID(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/user/%d/", users[2].ID), nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res map[string]string

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["error"] != "Access denied" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Access denied'. Got '%s'",
			res["error"],
		)
	}
}

func TestDeleteUserNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/user/777/", nil)

	db := models.InitTestDbURI(
		"postgresql://root@localhost:26257/todogo_test?sslmode=disable", false,
	)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}
