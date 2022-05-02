package main

import (
	"fmt"
	"log"
	"net/http"

	"goswagger/app/students"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/students", students.GetStudents).Methods("GET")
	r.HandleFunc("/students/{id}", students.GetStudent).Methods("GET")
	r.HandleFunc("/students", students.CreateStudent).Methods("POST")

	// starting the server
	fmt.Printf("Starting the server at 8000...\n")

	log.Fatal(http.ListenAndServe(":8000", r))
}
