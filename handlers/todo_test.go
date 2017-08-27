package handlers

import (
	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/models"
	"testing"
	"net/http"
	"encoding/json"
	"fmt"
	"bytes"
	"sort"
)

var shoppingTodos = []models.Todo{
	{ID: 1, Title: "Milk", Completed: true, Note: "Milk", TodoListID: 1, UserID: users[0].ID},
	{ID: 2, Title: "Bread", Completed: false, Note: "Bread", TodoListID: 1, UserID: users[0].ID},
	{ID: 3, Title: "Cucumber", Completed: true, Note: "Cucumber", TodoListID: 1, UserID: users[0].ID},
	{ID: 4, Title: "Tomato", Completed: false, Note: "Tomato", TodoListID: 1, UserID: users[0].ID},
	{ID: 5, Title: "Oil", Completed: false, Note: "Oil", TodoListID: 1, UserID: users[0].ID},
	{ID: 6, Title: "Potato", Completed: false, Note: "Potato", TodoListID: 1, UserID: users[0].ID},
	{ID: 7, Title: "Ice cream", Completed: true, Note: "Ice cream", TodoListID: 1, UserID: users[0].ID},
}
var workTodos = []models.Todo{
	{ID: 8, Title: "Write some tests for todo list", Completed: true, Note: "", TodoListID: 2, UserID: users[0].ID},
	{ID: 9, Title: "Write some tests for todo", Completed: false, Note: "", TodoListID: 2, UserID: users[0].ID},
	{ID: 10, Title: "Implement authentication", Completed: false, Note: "", TodoListID: 2, UserID: users[0].ID},
	{ID: 11, Title: "Implement frontend in clojure script", Completed: false, Note: "", TodoListID: 2, UserID: users[0].ID},
}

func CreateTodoFixtures(db *gorm.DB) {
	for _, todo := range shoppingTodos {
		db.Create(&todo)
	}
	for _, todo := range workTodos {
		db.Create(&todo)
	}
}

func TestGetTodos(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/1/todo/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res []models.Todo

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	// sort TODOs to display incomplete tasks first
	sort.Slice(shoppingTodos, func(i, j int) bool {
		var completedI, completedJ int8
		if shoppingTodos[i].Completed {
			completedI = 1
		}
		if shoppingTodos[j].Completed {
			completedJ = 1
		}

		if completedI < completedJ {
			return true
		}
		if completedI > completedJ {
			return false
		}

		return shoppingTodos[i].ID < shoppingTodos[j].ID
	})

	for i := range shoppingTodos {
		if res[i].ID != shoppingTodos[i].ID {
			t.Errorf("Response body should be `%s`, was  %s", shoppingTodos[i].ID, res[i].ID)
		}
		if res[i].Title != shoppingTodos[i].Title {
			t.Errorf("Response body should be `%s`, was  %s", shoppingTodos[i].Title, res[i].Title)
		}
		if res[i].Completed != shoppingTodos[i].Completed {
			t.Errorf("Response body should be `%s`, was  %s", shoppingTodos[i].Completed, res[i].Completed)
		}
		if res[i].Note != shoppingTodos[i].Note {
			t.Errorf("Response body should be `%s`, was  %s", shoppingTodos[i].Note, res[i].Note)
		}
		if res[i].TodoListID != shoppingTodos[i].TodoListID {
			t.Errorf("Response body should be `%s`, was  %s", shoppingTodos[i].TodoListID, res[i].TodoListID)
		}
	}
}

func TestGetTodosNoTodos(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/3/todo/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res []models.Todo

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(res) != 0 {
		t.Errorf("Response body should be `[]`, was  %s", res)
	}
}

func TestGetTodosNonExistentList(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/777/todo/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

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

func TestGetTodosNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/1/todo/", nil)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestGetTodo(t *testing.T) {
	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"GET", fmt.Sprintf("/api/v1/list/%d/todo/%d/", todo.TodoListID, todo.ID), nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res models.Todo

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.Title != todo.Title {
		t.Errorf("Response body should be `%s`, was  %s", todo.Title, res.Title)
	}
	if res.Completed != todo.Completed {
		t.Errorf("Response body should be `%s`, was  %s", todo.Completed, res.Completed)
	}
	if res.Note != todo.Note {
		t.Errorf("Response body should be `%s`, was  %s", todo.Note, res.Note)
	}
	if res.TodoListID != todo.TodoListID {
		t.Errorf("Response body should be `%s`, was  %s", todo.TodoListID, res.TodoListID)
	}
}

func TestGetTodoWrongListID(t *testing.T) {
	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"GET", fmt.Sprintf("/api/v1/list/%d/todo/%d/", 3, todo.ID), nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

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

	if res["error"] != "Todo not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo List not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestGetTodoWrongID(t *testing.T) {
	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"GET", fmt.Sprintf("/api/v1/list/%d/todo/%d/", todo.TodoListID, 777), nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

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

	if res["error"] != "Todo not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo List not found'. Got '%s'",
			res["error"],
		)
	}
}


func TestGetTodoNoAuthToken(t *testing.T) {
	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"GET", fmt.Sprintf("/api/v1/list/%d/todo/%d/", todo.TodoListID, todo.ID), nil,
	)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestCreateTodo(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	req, _ := http.NewRequest(
		"POST", "/api/v1/list/1/todo/", bytes.NewBuffer(jsonStr),
	)
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

	var res models.Todo

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.Title != "Milk" {
		t.Errorf("Response body should be `Milk`, was %s", res.Title)
	}
	if !res.Completed {
		t.Errorf("Response body should be `true`, was %s", res.Completed)
	}
	if res.Note != "1.5 L 1.5%" {
		t.Errorf("Response body should be `1.5 L 1.5%%`, was %s", res.Note)
	}
	if res.TodoListID != 1 {
		t.Errorf("Response body should be `1`, was %s", res.TodoListID)
	}
}

func TestCreateTodoNonExistentList(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	req, _ := http.NewRequest(
		"POST", "/api/v1/list/777/todo/", bytes.NewBuffer(jsonStr),
	)
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

func TestCreateTodoNoAuthToken(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	req, _ := http.NewRequest(
		"POST", "/api/v1/list/1/todo/", bytes.NewBuffer(jsonStr),
	)

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

func TestUpdateTodo(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("/api/v1/list/%d/todo/%d/", todo.TodoListID, todo.ID),
		bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString := resp.Body.String()

	var res models.Todo

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.Title != "Milk" {
		t.Errorf("Response body should be `Milk`, was %s", res.Title)
	}
	if !res.Completed {
		t.Errorf("Response body should be `true`, was %s", res.Completed)
	}
	if res.Note != "1.5 L 1.5%" {
		t.Errorf("Response body should be `1.5 L 1.5%%`, was %s", res.Note)
	}
	if res.TodoListID != 1 {
		t.Errorf("Response body should be `1`, was %s", res.TodoListID)
	}
}

func TestUpdateTodoWrongListID(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("/api/v1/list/%d/todo/%d/", 777, todo.ID),
		bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

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

func TestUpdateTodoWrongID(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("/api/v1/list/%d/todo/%d/", todo.TodoListID, 777),
		bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

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

	if res["error"] != "Todo not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestUpdateTodoNoAuthToken(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("/api/v1/list/%d/todo/%d/", todo.TodoListID, todo.ID),
		bytes.NewBuffer(jsonStr),
	)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}

func TestUpdateTodoWrongID1(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("/api/v1/list/%d/todo/%d/", 3, todo.ID),
		bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

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

	if res["error"] != "Todo not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestDeleteTodo(t *testing.T) {
	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/list/%d/todo/%d/", todo.TodoListID, todo.ID),
		nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusNoContent, resp.Code)
	}
}

func TestDeleteTodoWrongListID(t *testing.T) {
	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/list/%d/todo/%d/", 777, todo.ID),
		nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

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

func TestDeleteTodoWrongID(t *testing.T) {
	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/list/%d/todo/%d/", 3, todo.ID),
		nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

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

	if res["error"] != "Todo not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestDeleteTodoNoAuthToken(t *testing.T) {
	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/list/%d/todo/%d/", todo.TodoListID, todo.ID),
		nil,
	)

	db := models.InitTestDb("localhost", "todogo", "todogo", "todogo", false)

	models.DropTables(db)
	models.CreateTables(db)
	CreateUserFixtures(db)
	CreateTodoListFixtures(db)
	CreateTodoFixtures(db)

	resp := ExecuteRequest(db, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, resp.Code)
	}
}
