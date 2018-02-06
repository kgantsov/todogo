package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kgantsov/todogo/models"
	uuid "github.com/satori/go.uuid"
)

var shoppingTodos = []models.Todo{
	{
		ID:         uuid.NewV4(),
		Title:      "Milk",
		Completed:  true,
		Note:       "Milk",
		TodoListID: todoLists[0].ID,
		UserID:     users[0].ID,
		Priority:   models.PRIORITY_NORMAL,
		CreatedAt:  ptr(time.Now().Add(time.Second * 2)),
	},
	{
		ID:         uuid.NewV4(),
		Title:      "Bread",
		Completed:  false,
		Note:       "Bread",
		TodoListID: todoLists[0].ID,
		UserID:     users[0].ID,
		Priority:   models.PRIORITY_NORMAL,
		CreatedAt:  ptr(time.Now().Add(time.Second * 3)),
	},
	{
		ID:         uuid.NewV4(),
		Title:      "Cucumber",
		Completed:  true,
		Note:       "Cucumber",
		TodoListID: todoLists[0].ID,
		UserID:     users[0].ID,
		Priority:   models.PRIORITY_NORMAL,
		CreatedAt:  ptr(time.Now().Add(time.Second * 4)),
	},
	{
		ID:         uuid.NewV4(),
		Title:      "Tomato",
		Completed:  false,
		Note:       "Tomato",
		TodoListID: todoLists[0].ID,
		UserID:     users[0].ID,
		Priority:   models.PRIORITY_URGENT,
		CreatedAt:  ptr(time.Now().Add(time.Second * 5)),
	},
	{
		ID:         uuid.NewV4(),
		Title:      "Oil",
		Completed:  false,
		Note:       "Oil",
		TodoListID: todoLists[0].ID,
		UserID:     users[0].ID,
		Priority:   models.PRIORITY_HIGH,
		CreatedAt:  ptr(time.Now().Add(time.Second * 6)),
	},
	{
		ID:         uuid.NewV4(),
		Title:      "Potato",
		Completed:  false,
		Note:       "Potato",
		TodoListID: todoLists[0].ID,
		UserID:     users[0].ID,
		Priority:   models.PRIORITY_NORMAL,
		CreatedAt:  ptr(time.Now().Add(time.Second * 7)),
	},
	{
		ID:         uuid.NewV4(),
		Title:      "Ice cream",
		Completed:  true,
		Note:       "Ice cream",
		TodoListID: todoLists[0].ID,
		UserID:     users[0].ID,
		Priority:   models.PRIORITY_NORMAL,
		CreatedAt:  ptr(time.Now().Add(time.Second * 8)),
	},
}

var workTodos = []models.Todo{
	{
		ID:         uuid.NewV4(),
		Title:      "Write some tests for todo list",
		Completed:  true,
		Note:       "",
		TodoListID: todoLists[1].ID,
		UserID:     users[0].ID,
		CreatedAt:  ptr(time.Now().Add(time.Second * 9)),
	},
	{
		ID:         uuid.NewV4(),
		Title:      "Write some tests for todo",
		Completed:  false,
		Note:       "",
		TodoListID: todoLists[1].ID,
		UserID:     users[0].ID,
		CreatedAt:  ptr(time.Now().Add(time.Second * 10)),
	},
	{
		ID:         uuid.NewV4(),
		Title:      "Implement authentication",
		Completed:  false,
		Note:       "",
		TodoListID: todoLists[1].ID,
		UserID:     users[0].ID,
		CreatedAt:  ptr(time.Now().Add(time.Second * 11)),
	},
	{
		ID:         uuid.NewV4(),
		Title:      "Implement frontend in clojure script",
		Completed:  false,
		Note:       "",
		TodoListID: todoLists[1].ID,
		UserID:     users[0].ID,
		CreatedAt:  ptr(time.Now().Add(time.Second * 12)),
	},
}

func ptr(t time.Time) *time.Time {
	return &t
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
	req, _ := http.NewRequest(
		"GET", fmt.Sprintf("/api/v1/list/%s/todo/", todoLists[0].ID.String()), nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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

		if shoppingTodos[i].Priority < shoppingTodos[j].Priority {
			return false
		}
		if shoppingTodos[i].Priority > shoppingTodos[j].Priority {
			return true
		}

		return shoppingTodos[i].CreatedAt.Before(*shoppingTodos[j].CreatedAt)
	})

	for i := range shoppingTodos {
		if res[i].ID != shoppingTodos[i].ID {
			t.Errorf("Response body should be `%v`, was  %v", shoppingTodos[i].ID, res[i].ID)
		}
		if res[i].Title != shoppingTodos[i].Title {
			t.Errorf("Response body should be `%s`, was  %s", shoppingTodos[i].Title, res[i].Title)
		}
		if res[i].Completed != shoppingTodos[i].Completed {
			t.Errorf("Response body should be `%v`, was  %v", shoppingTodos[i].Completed, res[i].Completed)
		}
		if res[i].Note != shoppingTodos[i].Note {
			t.Errorf("Response body should be `%s`, was  %s", shoppingTodos[i].Note, res[i].Note)
		}
		if res[i].TodoListID != shoppingTodos[i].TodoListID {
			t.Errorf("Response body should be `%v`, was  %v", shoppingTodos[i].TodoListID, res[i].TodoListID)
		}
		if res[i].UserID != shoppingTodos[i].UserID {
			t.Errorf("Response body should be `%v`, was  %v", shoppingTodos[i].UserID, res[i].UserID)
		}
	}
}

func TestGetTodosNoTodos(t *testing.T) {
	req, _ := http.NewRequest(
		"GET", fmt.Sprintf("/api/v1/list/%s/todo/", todoLists[2].ID.String()), nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if len(res) != 0 {
		t.Errorf("Response body should be `[]`, was %v", res)
	}
}

func TestGetTodosNonExistentList(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/list/E41B72FE-B184-4A85-B280-0544DEC106C7/todo/", nil)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if res["error"] != "Todo list not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo list not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestGetTodosNoAuthToken(t *testing.T) {
	req, _ := http.NewRequest(
		"GET", fmt.Sprintf("/api/v1/list/%s/todo/", todoLists[2].ID.String()), nil,
	)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
		"GET",
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todo.TodoListID.String(), todo.ID.String()),
		nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if res.Title != todo.Title {
		t.Errorf("Response body should be `%s`, was  %s", todo.Title, res.Title)
	}
	if res.Completed != todo.Completed {
		t.Errorf("Response body should be `%v`, was  %v", todo.Completed, res.Completed)
	}
	if res.Note != todo.Note {
		t.Errorf("Response body should be `%s`, was  %s", todo.Note, res.Note)
	}
	if res.TodoListID != todo.TodoListID {
		t.Errorf("Response body should be `%s`, was  %s", todo.TodoListID.String(), res.TodoListID.String())
	}
	if res.UserID != todo.UserID {
		t.Errorf("Response body should be `%s`, was  %s", todo.UserID.String(), res.UserID.String())
	}
}

func TestGetTodoWrongListID(t *testing.T) {
	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todoLists[2].ID.String(), todo.ID.String()),
		nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
		"GET",
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todo.TodoListID.String(), "BC7A315B-E463-4436-A11A-E6D902F05214"),
		nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
		"GET",
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todo.TodoListID.String(), todo.ID.String()),
		nil,
	)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
		"POST",
		fmt.Sprintf("/api/v1/list/%s/todo/", todoLists[0].ID.String()),
		bytes.NewBuffer(jsonStr),
	)
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

	var todo models.Todo

	err := json.Unmarshal([]byte(bodyAsString), &todo)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	if todo.Title != "Milk" {
		t.Errorf("Response body should be `Milk`, was %s", todo.Title)
	}
	if !todo.Completed {
		t.Errorf("Response body should be `true`, was %v", todo.Completed)
	}
	if todo.Note != "1.5 L 1.5%" {
		t.Errorf("Response body should be `1.5 L 1.5%%`, was %s", todo.Note)
	}
	if todo.TodoListID != todoLists[0].ID {
		t.Errorf("Response body should be `%v`, was %v", todoLists[0].ID, todo.TodoListID)
	}
	if todo.Priority != 4 {
		t.Errorf("Response body should be `4`, was %d", todo.Priority)
	}

	req, _ = http.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todo.TodoListID.String(), todo.ID.String()),
		nil,
	)
	req.Header.Set("Auth-Token", token)

	resp = ExecuteRequest(db, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	bodyAsString = resp.Body.String()

	var res models.Todo

	err = json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
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
		t.Errorf("Response body should be `%d`, was  %d", todo.TodoListID, res.TodoListID)
	}
	if res.UserID != todo.UserID {
		t.Errorf("Response body should be `%d`, was  %d", todo.UserID, res.UserID)
	}
}

func TestCreateTodoWithPriority(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%", "priority": 5}`)

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("/api/v1/list/%s/todo/", todoLists[0].ID.String()),
		bytes.NewBuffer(jsonStr),
	)
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

	var res models.Todo

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	if res.Title != "Milk" {
		t.Errorf("Response body should be `Milk`, was %s", res.Title)
	}
	if !res.Completed {
		t.Errorf("Response body should be `true`, was %v", res.Completed)
	}
	if res.Note != "1.5 L 1.5%" {
		t.Errorf("Response body should be `1.5 L 1.5%%`, was %s", res.Note)
	}
	if res.TodoListID != todoLists[0].ID {
		t.Errorf("Response body should be `%v`, was %v", todoLists[0].ID, res.TodoListID)
	}
	if res.Priority != 5 {
		t.Errorf("Response body should be `5`, was %d", res.Priority)
	}
}

func TestCreateTodoWithDeadLine(t *testing.T) {
	var jsonStr = []byte(
		`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%", "dead_line_at": "2009-11-17T20:34:58.0Z"}`,
	)

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("/api/v1/list/%s/todo/", todoLists[0].ID.String()),
		bytes.NewBuffer(jsonStr),
	)
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

	var res models.Todo

	err := json.Unmarshal([]byte(bodyAsString), &res)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	if res.Title != "Milk" {
		t.Errorf("Response body should be `Milk`, was %s", res.Title)
	}
	if !res.Completed {
		t.Errorf("Response body should be `true`, was %v", res.Completed)
	}
	if res.Note != "1.5 L 1.5%" {
		t.Errorf("Response body should be `1.5 L 1.5%%`, was %s", res.Note)
	}
	if res.TodoListID != todoLists[0].ID {
		t.Errorf("Response body should be `%v`, was %v", todoLists[0].ID, res.TodoListID)
	}
	if *res.DeadLineAt != time.Date(2009, 11, 17, 20, 34, 58, 0, time.UTC) {
		t.Errorf("Response body should be `1`, was %s", res.DeadLineAt)
	}
}

func TestCreateTodoNonExistentList(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	req, _ := http.NewRequest(
		"POST", "/api/v1/list/BC7A315B-E463-4436-A11A-E6D902F05214/todo/", bytes.NewBuffer(jsonStr),
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

	if res["error"] != "Todo list not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo list not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestCreateTodoNoAuthToken(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("/api/v1/list/%s/todo/", todoLists[0].ID.String()),
		bytes.NewBuffer(jsonStr),
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

func TestUpdateTodo(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todo.TodoListID.String(), todo.ID.String()),
		bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if res.Title != "Milk" {
		t.Errorf("Response body should be `Milk`, was %s", res.Title)
	}
	if !res.Completed {
		t.Errorf("Response body should be `true`, was %v", res.Completed)
	}
	if res.Note != "1.5 L 1.5%" {
		t.Errorf("Response body should be `1.5 L 1.5%%`, was %s", res.Note)
	}
	if res.TodoListID != todo.TodoListID {
		t.Errorf("Response body should be `%v`, was %v", todo.TodoListID, res.TodoListID)
	}
}

func TestUpdateTodoWithPriority(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%", "priority": 7}`)

	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todo.TodoListID.String(), todo.ID.String()),
		bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if res.Title != "Milk" {
		t.Errorf("Response body should be `Milk`, was %s", res.Title)
	}
	if !res.Completed {
		t.Errorf("Response body should be `true`, was %v", res.Completed)
	}
	if res.Note != "1.5 L 1.5%" {
		t.Errorf("Response body should be `1.5 L 1.5%%`, was %s", res.Note)
	}
	if res.TodoListID != todo.TodoListID {
		t.Errorf("Response body should be `%v`, was %v", todo.TodoListID, res.TodoListID)
	}
	if res.Priority != 7 {
		t.Errorf("Response body should be `7`, was %d", res.Priority)
	}
}

func TestUpdateTodoWrongListID(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("/api/v1/list/BC7A315B-E463-4436-A11A-E6D902F05214/todo/%s/", todo.ID.String()),
		bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if res["error"] != "Todo list not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo list not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestUpdateTodoWrongID(t *testing.T) {
	var jsonStr = []byte(`{"title": "Milk", "completed": true, "note": "1.5 L 1.5%"}`)

	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"PUT",
		fmt.Sprintf("/api/v1/list/%s/todo/BC7A315B-E463-4436-A11A-E6D902F05214/", todo.TodoListID.String()),
		bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todo.TodoListID.String(), todo.ID.String()),
		bytes.NewBuffer(jsonStr),
	)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todoLists[2].ID.String(), todo.ID.String()),
		bytes.NewBuffer(jsonStr),
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todo.TodoListID.String(), todo.ID.String()),
		nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
		fmt.Sprintf("/api/v1/list/BC7A315B-E463-4436-A11A-E6D902F05214/todo/%s/", todo.ID.String()),
		nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
	}

	if res["error"] != "Todo list not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'Todo list not found'. Got '%s'",
			res["error"],
		)
	}
}

func TestDeleteTodoWrongID(t *testing.T) {
	todo := shoppingTodos[2]

	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todoLists[2].ID.String(), todo.ID.String()),
		nil,
	)
	token, _ := createToken(users[0].ID)
	req.Header.Set("Auth-Token", token)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
		fmt.Sprintf("/api/v1/list/%s/todo/%s/", todo.TodoListID.String(), todo.ID.String()),
		nil,
	)

	db := models.InitTestDbURI(testDBConnectionString, false)

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
