/*****
Title: Operata Technical Assessment - Unit Test go rest api server
Author: Anitha Kurian
******/
package main

import(
"testing"
"net/http"
"net/http/httptest"
"bytes"
)

func TestCreateGetUserGetUsers(t *testing.T) {
	var jsonStr = []byte(`{"FirstName":"Anitha","LastName":"Kurian","ID":"ID1","Age":31}`)
	reqCreateUser, errCreateUser := http.NewRequest("POST", "localhost:8080/users", bytes.NewBuffer(jsonStr))
	if errCreateUser != nil {
		t.Fatal(errCreateUser)
	}
	reqCreateUser.Header.Set("Content-Type", "application/json")
	rrCreateUser := httptest.NewRecorder()
	userHandlersCreateUser := newUserHandlers()
	handlerCreateUser := http.HandlerFunc(userHandlersCreateUser.users)
	handlerCreateUser.ServeHTTP(rrCreateUser, reqCreateUser)
	if statusCreateUser := rrCreateUser.Code; statusCreateUser != http.StatusOK {
		t.Errorf("handlerCreateUser returned wrong statusCreateUser code: got %v want %v",
		statusCreateUser, http.StatusOK)
	}
	expectedCreateUser := `{"FirstName":"Anitha","LastName":"Kurian","ID":"ID1","Age":31}`
	if rrCreateUser.Body.String() != expectedCreateUser {
		t.Errorf("handlerCreateUser returned unexpected body: got %v want %v",
		rrCreateUser.Body.String(), expectedCreateUser)
	}

	reqGetUsers, errGetUsers := http.NewRequest("GET", "localhost:8080/users", nil)
	if errGetUsers != nil {
		t.Fatal(errGetUsers)
	}
    rrGetUsers := httptest.NewRecorder()
	handlerGetUsers := http.HandlerFunc(userHandlersCreateUser.users)
	handlerGetUsers.ServeHTTP(rrGetUsers, reqGetUsers)
	if statusGetUsers := rrGetUsers.Code; statusGetUsers != http.StatusOK{
		t.Errorf("handlerGetUsers returned wrong status code: got %v want %v",
		statusGetUsers, http.StatusOK)}

	// Check the response body
	expectedGetUsers := `[{"FirstName":"Anitha","LastName":"Kurian","ID":"ID1","Age":31}]`
    if rrGetUsers.Body.String() != expectedGetUsers{
		t.Errorf("handlerGetUsers returned unexpected body: got %v want %v",
		rrGetUsers.Body.String(), expectedGetUsers)
	}

	reqGetUser, errGetUser := http.NewRequest("GET", "/users/ID1", nil)
	if errGetUser != nil {
		t.Fatal(errGetUser)
	}
	
	rrGetUser := httptest.NewRecorder()
	handlerGetUser := http.HandlerFunc(userHandlersCreateUser.getUser)
	handlerGetUser.ServeHTTP(rrGetUser, reqGetUser)
	if status := rrGetUser.Code; status != http.StatusOK {
		t.Errorf("handlerGetUser returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expectedGetUser := `{"FirstName":"Anitha","LastName":"Kurian","ID":"ID1","Age":31}`
	if rrGetUser.Body.String() != expectedGetUser {
		t.Errorf("handlerGetUser returned unexpected body: got %v want %v",
		rrGetUser.Body.String(), expectedGetUser)
	}

}

