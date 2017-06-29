package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/models"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func DBMiddleware(db gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("db", db)

		c.Next()
	}
}

func CreateFixtures(db *gorm.DB) {
	db.Create(&models.TodoList{Title: "Shopping list"})
	db.Create(&models.TodoList{Title: "Work list"})
	db.Create(&models.TodoList{Title: "Sport list"})
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(db *gorm.DB, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(DBMiddleware(*db))

	DefineRoutes(r)

	r.ServeHTTP(w, req)

	return w
}

func TestGetTodoLists(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/", nil)

	db := models.InitTestDb()

	models.DropTables(db)
	models.CreateTables(db)
	CreateFixtures(db)

	resp := executeRequest(db, req)

	checkResponseCode(t, http.StatusOK, resp.Code)

	bodyAsString := resp.Body.String()

	var res []models.TodoList

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
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

	db := models.InitTestDb()

	models.DropTables(db)
	models.CreateTables(db)

	resp := executeRequest(db, req)

	checkResponseCode(t, http.StatusOK, resp.Code)

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

func TestGetTodoList(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/2/", nil)

	db := models.InitTestDb()

	models.DropTables(db)
	models.CreateTables(db)
	CreateFixtures(db)

	resp := executeRequest(db, req)

	checkResponseCode(t, http.StatusOK, resp.Code)

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

	db := models.InitTestDb()

	models.DropTables(db)
	models.CreateTables(db)
	CreateFixtures(db)

	resp := executeRequest(db, req)

	checkResponseCode(t, http.StatusNotFound, resp.Code)

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

func TestCreateTodoList(t *testing.T) {
	var jsonStr = []byte(`{"title": "My tasks"}`)

	req, _ := http.NewRequest("POST", "/api/v1/list/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDb()

	models.DropTables(db)
	models.CreateTables(db)
	CreateFixtures(db)

	resp := executeRequest(db, req)

	checkResponseCode(t, http.StatusCreated, resp.Code)

	bodyAsString := resp.Body.String()

	var res models.TodoList

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.ID != 4 {
		t.Errorf("Response body should be `4`, was  %s", res.ID)
	}
	if res.Title != "My tasks" {
		t.Errorf("Response body should be `My tasks`, was  %s", res.Title)
	}
}

func TestCreateTodoListMissedFields(t *testing.T) {
	var jsonStr = []byte(`{}`)

	req, _ := http.NewRequest("POST", "/api/v1/list/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDb()

	models.DropTables(db)
	models.CreateTables(db)
	CreateFixtures(db)

	resp := executeRequest(db, req)

	checkResponseCode(t, http.StatusBadRequest, resp.Code)
}

func TestUpdateTodoList(t *testing.T) {
	var jsonStr = []byte(`{"title": "My tasks"}`)

	req, _ := http.NewRequest("PUT", "/api/v1/list/2/", bytes.NewBuffer(jsonStr))

	db := models.InitTestDb()

	models.DropTables(db)
	models.CreateTables(db)
	CreateFixtures(db)

	resp := executeRequest(db, req)

	checkResponseCode(t, http.StatusOK, resp.Code)

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

	db := models.InitTestDb()

	models.DropTables(db)
	models.CreateTables(db)
	CreateFixtures(db)

	resp := executeRequest(db, req)

	checkResponseCode(t, http.StatusNotFound, resp.Code)

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

func TestDeleteTodoList(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/list/2/", nil)

	db := models.InitTestDb()

	models.DropTables(db)
	models.CreateTables(db)
	CreateFixtures(db)

	resp := executeRequest(db, req)

	checkResponseCode(t, http.StatusNoContent, resp.Code)
}

func TestDeleteTodoListWrongID(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/list/777/", nil)

	db := models.InitTestDb()

	models.DropTables(db)
	models.CreateTables(db)
	CreateFixtures(db)

	resp := executeRequest(db, req)

	checkResponseCode(t, http.StatusNotFound, resp.Code)

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
