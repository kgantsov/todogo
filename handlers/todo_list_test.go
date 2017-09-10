package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/models"
	"gopkg.in/gin-gonic/gin.v1"
)

func DBMiddleware(db gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("db", db)

		c.Next()
	}
}

func CreateTodoListFixtures(db *gorm.DB) {
	db.Create(&models.TodoList{ID: 1, Title: "Shopping list", UserID: users[0].ID})
	db.Create(&models.TodoList{ID: 2, Title: "Work list", UserID: users[0].ID})
	db.Create(&models.TodoList{ID: 3, Title: "Sport list", UserID: users[0].ID})
	db.Create(&models.TodoList{ID: 4, Title: "Todo project", UserID: users[1].ID})
	db.Create(&models.TodoList{ID: 5, Title: "Foogle project", UserID: users[1].ID})
	db.Create(&models.TodoList{ID: 6, Title: "Sport list", UserID: users[2].ID})
}

func ExecuteRequest(db *gorm.DB, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(DBMiddleware(*db))

	DefineRoutes(db, r)

	r.ServeHTTP(w, req)

	return w
}

func TestGetTodoLists(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res []models.TodoList

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(res) != 3 {
		t.Errorf("Response body should be `3`, was  %s", len(res))
	}

	if res[0].ID != 1 {
		t.Errorf("Response body should be `1`, was  %s", res[0].ID)
	}
	if res[0].Title != "Shopping list" {
		t.Errorf("Response body should be `Shopping list`, was  %s", res[0].Title)
	}

	if res[1].ID != 2 {
		t.Errorf("Response body should be `2`, was  %s", res[1].ID)
	}
	if res[1].Title != "Work list" {
		t.Errorf("Response body should be `Work list`, was  %s", res[1].Title)
	}

	if res[2].ID != 3 {
		t.Errorf("Response body should be `3`, was  %s", res[2].ID)
	}
	if res[2].Title != "Sport list" {
		t.Errorf("Response body should be `Sport list`, was  %s", res[2].Title)
	}
}

func TestGetTodoListsEmptyTable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/", nil)
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

	var res []models.TodoList

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(res) != 0 {
		t.Errorf("Response body should be empty, was  %v", res)
	}
}

func TestGetTodoListsNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/", nil)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestGetTodoListsAuthUserDoesNotExist(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/", nil)
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

func TestGetTodoListsWrongAuthToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/", nil)
	req.Header.Set("Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDE0NDg3MDr")

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestGetTodoList(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/2/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res models.TodoList

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.ID != 2 {
		t.Errorf("Response body should be `2`, was  %s", res.ID)
	}
	if res.Title != "Work list" {
		t.Errorf("Response body should be `Work list`, was  %s", res.Title)
	}
}

func TestGetTodoListWrongID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/777/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

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

	if res["error"] != "TODO list not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'TODO list not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestGetTodoListNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/2/", nil)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestCreateTodoList(t *testing.T) {
	var jsonStr = []byte(`{"title": "My tasks"}`)

	req, _ := http.NewRequest("POST", "/api/v1/list/", bytes.NewBuffer(jsonStr))
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusCreated, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res models.TodoList

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.Title != "My tasks" {
		t.Errorf("Response body should be `My tasks`, was  %s", res.Title)
	}
}

func TestCreateTodoListMissedFields(t *testing.T) {
	var jsonStr = []byte(`{}`)

	req, _ := http.NewRequest("POST", "/api/v1/list/", bytes.NewBuffer(jsonStr))
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusBadRequest, resp.Code)
	}
}

func TestCreateTodoListNoAuthToken(t *testing.T) {
	var jsonStr = []byte(`{}`)

	req, _ := http.NewRequest("POST", "/api/v1/list/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestUpdateTodoList(t *testing.T) {
	var jsonStr = []byte(`{"title": "My tasks"}`)

	req, _ := http.NewRequest("PUT", "/api/v1/list/2/", bytes.NewBuffer(jsonStr))
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res models.TodoList

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.ID != 2 {
		t.Errorf("Response body should be `4`, was  %s", res.ID)
	}
	if res.Title != "My tasks" {
		t.Errorf("Response body should be `My tasks`, was  %s", res.Title)
	}
}

func TestUpdateTodoListWrongID(t *testing.T) {
	var jsonStr = []byte(`{"title": "My tasks"}`)

	req, _ := http.NewRequest("PUT", "/api/v1/list/777/", bytes.NewBuffer(jsonStr))
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

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

	if res["error"] != "Todo List not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo List not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestUpdateTodoListNoAuthToken(t *testing.T) {
	var jsonStr = []byte(`{"title": "My tasks"}`)

	req, _ := http.NewRequest("PUT", "/api/v1/list/777/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestDeleteTodoList(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/list/2/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusNoContent, resp.Code)
	}
}

func TestDeleteTodoListWrongID(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/list/777/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

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

	if res["error"] != "Todo List not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo List not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestDeleteTodoListNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/list/777/", nil)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}
