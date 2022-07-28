// Package classification People API.
//
// Documentation for Person API
//
//     Schemes: http
//     BasePath: /
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
// 	swagger:meta
package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

type Person struct {
	ID        string   `json:"id,omitempty"`
	FirstName string   `json:"firstname,omitempty"`
	LastName  string   `json:"lastname,omitempty"`
	Addres    *Address `json:"addres,omitempty"`
}

var people []Person

// A list of persons returns in the responses
// swagger:response personResponse
type peopleResponse struct {
	// All people in the listsystem
	// in: body
	Body []Person
}

// No content
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:parameters blablabla DeletePerson
type peopleIDParameterWraper struct {
	// The id of the person to delete
	// in: path
	// required: true
	ID int `json:"id"`
}

func GetInput(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
}

// swagger:route GET /people Person GetPerson
// Returns a list of person
// responses:
//    200: personResponse

// GetPerson returns the person
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// swagger:route DELETE /people/{id} Person DeletePerson
// Delete person from the database
//
// responses:
//	201: noContentResponse

// DeletePerson delete a person from the database
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(people)
}

func main() {
	router := mux.NewRouter()

	people = append(people, Person{ID: "1", FirstName: "John", LastName: "Doe", Addres: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", FirstName: "Petter", LastName: "Tomson"})

	//endpoints
	router.HandleFunc("/people", GetInput).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

	options := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(options, nil)

	router.Handle("/docs", sh)
	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	http.ListenAndServe(":3333", router)

}
