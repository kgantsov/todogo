package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kgantsov/todogo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func DBMiddleware(db gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("db", db)

		c.Next()
	}
}

var todoLists = []models.TodoList{
	{
		ID:        uuid.NewV4(),
		Title:     "Shopping list",
		UserID:    users[0].ID,
		CreatedAt: ptr(time.Now().Add(time.Second * 2)),
	},
	{
		ID:        uuid.NewV4(),
		Title:     "Work list",
		UserID:    users[0].ID,
		CreatedAt: ptr(time.Now().Add(time.Second * 3)),
	},
	{
		ID:        uuid.NewV4(),
		Title:     "Sport list",
		UserID:    users[0].ID,
		CreatedAt: ptr(time.Now().Add(time.Second * 4)),
	},
	{
		ID:        uuid.NewV4(),
		Title:     "Todo project",
		UserID:    users[1].ID,
		CreatedAt: ptr(time.Now().Add(time.Second * 5)),
	},
	{
		ID:        uuid.NewV4(),
		Title:     "Foogle project",
		UserID:    users[1].ID,
		CreatedAt: ptr(time.Now().Add(time.Second * 6)),
	},
	{
		ID:        uuid.NewV4(),
		Title:     "Sport list",
		UserID:    users[2].ID,
		CreatedAt: ptr(time.Now().Add(time.Second * 7)),
	},
}

func CreateTodoListFixtures(db *gorm.DB) {
	for _, todoList := range todoLists {
		db.Create(&todoList)
	}
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

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if len(res) != 3 {
		t.Errorf("Response body should be `3`, was  %d", len(res))
	}

	if res[0].ID != todoLists[0].ID {
		t.Errorf("Response body should be `%v`, was  %v", todoLists[0].ID, res[0].ID)
	}
	if res[0].Title != "Shopping list" {
		t.Errorf("Response body should be `Shopping list`, was  %s", res[0].Title)
	}

	if res[1].ID != todoLists[1].ID {
		t.Errorf("Response body should be `%v`, was  %v", todoLists[1].ID, res[1].ID)
	}
	if res[1].Title != "Work list" {
		t.Errorf("Response body should be `Work list`, was  %s", res[1].Title)
	}

	if res[2].ID != todoLists[2].ID {
		t.Errorf("Response body should be `%v`, was  %v", todoLists[2].ID, res[2].ID)
	}
	if res[2].Title != "Sport list" {
		t.Errorf("Response body should be `Sport list`, was  %s", res[2].Title)
	}
}

func TestGetTodoListsEmptyTable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if len(res) != 0 {
		t.Errorf("Response body should be empty, was  %v", res)
	}
}

func TestGetTodoListsNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/", nil)

	db := models.InitTestDbURI(testDBConnectionString, false)

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

	db := models.InitTestDbURI(testDBConnectionString, false)

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

	db := models.InitTestDbURI(testDBConnectionString, false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestGetTodoList(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/list/%s/", todoLists[2].ID.String()), nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if res.ID != todoLists[2].ID {
		t.Errorf("Response body should be `%v`, was  %v", todoLists[2].ID, res.ID)
	}
	if res.Title != todoLists[2].Title {
		t.Errorf("Response body should be `%s`, was  %s", todoLists[2].Title, res.Title)
	}
}

func TestGetTodoListWrongID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/E41B72FE-B184-4A85-B280-0544DEC106C7/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if res["error"] != "TODO list not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'TODO list not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestGetTodoListNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/E41B72FE-B184-4A85-B280-0544DEC106C7/", nil)

	db := models.InitTestDbURI(testDBConnectionString, false)

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

	var todoList models.TodoList

	err := json.Unmarshal([]byte(bodyAsString), &todoList)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	if todoList.Title != "My tasks" {
		t.Errorf("Response body should be `My tasks`, was  %s", todoList.Title)
	}

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/list/%s/", todoList.ID.String()), nil)
	req.Header.Set("Auth-Token", token)

	db = models.InitTestDbURI(os.Getenv("DB_TEST_URI"), false)

	resp = ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString = resp.Body.String()

	var res models.TodoList

	err = json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	if res.ID != todoList.ID {
		t.Errorf("Response body should be `%d`, was  %d", todoList.ID, res.ID)
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

	db := models.InitTestDbURI(testDBConnectionString, false)

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

	db := models.InitTestDbURI(testDBConnectionString, false)

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

	req, _ := http.NewRequest(
		"PUT", fmt.Sprintf("/api/v1/list/%s/", todoLists[2].ID.String()), bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if res.ID != todoLists[2].ID {
		t.Errorf("Response body should be `%v`, was  %v", todoLists[2].ID, res.ID)
	}
	if res.Title != "My tasks" {
		t.Errorf("Response body should be `My tasks`, was  %s", res.Title)
	}
}

func TestUpdateTodoListWrongID(t *testing.T) {
	var jsonStr = []byte(`{"title": "My tasks"}`)

	req, _ := http.NewRequest(
		"PUT", "/api/v1/list/E41B72FE-B184-4A85-B280-0544DEC106C7/", bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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

	req, _ := http.NewRequest(
		"PUT", "/api/v1/list/E41B72FE-B184-4A85-B280-0544DEC106C7/", bytes.NewBuffer(jsonStr),
	)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	req, _ := http.NewRequest(
		"DELETE", fmt.Sprintf("/api/v1/list/%s/", todoLists[2].ID.String()), nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	req, _ := http.NewRequest("DELETE", "/api/v1/list/E41B72FE-B184-4A85-B280-0544DEC106C7/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if res["error"] != "Todo List not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo List not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestDeleteTodoListNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/list/E41B72FE-B184-4A85-B280-0544DEC106C7/", nil)

	db := models.InitTestDbURI(testDBConnectionString, false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}
