/*****
Title: Operata Technical Assessment - Create go rest api server
Author: Anitha Kurian
******/
package main

import ( "net/http";
"encoding/json";
"sync";
"io/ioutil";
"fmt";
"strings"
 )

type User struct {
	FirstName string
	LastName string
	ID string
	Age int
}

// userHandlers is safe to use concurrently
type userHandlers struct{
	sync.Mutex
	store map[string] User
}

func (h *userHandlers) users(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case "GET":
		h.getUsers(w, r)
		return
	case "POST":
		h.createUser(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
	
}

//1. Create a user. A POST request with fields FirstName, LastName and Age
func (h *userHandlers) createUser(w http.ResponseWriter, r *http.Request){
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json"{
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json' but got '%s'", ct)))
		return
	}

	var user User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	h.Lock()
	// Lock so only one goroutine at a time can access the map.
	h.store[user.ID] = user
	defer h.Unlock()

	respondWithJSON(w, http.StatusOK, user)
		
}

// 2. Retrieve user by ID, A GET request.
func (h *userHandlers) getUser(w http.ResponseWriter, r *http.Request){

	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	user, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// 3. List all users. A GET request.
func (h *userHandlers) getUsers(w http.ResponseWriter, r *http.Request){
	users := make([]User, len(h.store))
	
	h.Lock()
	// Lock so only one goroutine at a time can access the map
	i := 0
	for _, user := range h.store {
		users[i] = user
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}


func newUserHandlers() *userHandlers{
	return &userHandlers{
		store: map[string]User{},
	}
}

// Called for responses to encode and send json data
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	//encode payload to json
	response, _ := json.Marshal(payload)

	//set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main(){
	userHandlers := newUserHandlers()
	http.HandleFunc("/users", userHandlers.users)
	http.HandleFunc("/users/", userHandlers.getUser)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil{
		panic(err)
	}
}