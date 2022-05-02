package students

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"goswagger/app/models"

	"github.com/gorilla/mux"
)

// swagger:model CommonSuccess
type CommonSuccess struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:model GetStud
type GetStuds struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string           `json:"message"`
	Data    []models.Student `json:"data"`
}

// swagger:model GetStud
type GetStud struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string `json:"message"`
	// Students for this user
	Data *models.Student `json:"data"`
}

// ErrHandler returns error message response
func ErrHandler(errmessage string) *models.CommonError {
	errresponse := models.CommonError{}
	errresponse.Status = 400
	errresponse.Message = errmessage
	return &errresponse
}

var students []models.Student

// swagger:route GET /students students listStudents
// Get students list
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// security:
// - apiKey: []
//
// responses:
//  401: CommonError
//  200: Student
func GetStudents(w http.ResponseWriter, r *http.Request) {
	response := GetStuds{}
	studentss := students
	//fmt.Println("he", studentss)

	response.Status = 1
	response.Message = "success"
	response.Data = studentss
	//fmt.Println(response.Data)

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// swagger:route  POST /students/{id} students findonestudent
// Find one student
//
// consumes:
//         - application/json
//
// responses:
//  401: CommonError
//  400: CommonError
//  200: Student
func GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range students {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(students)
}

// swagger:route POST /students students createStudent
// Create a new students
//
// consumes:
//         - application/json
//
// security:
// - apiKey: []
//
// responses:
//  401: CommonError
//  400: CommonError
//  200: Student
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var student models.Student
	// decode the body
	_ = json.NewDecoder(r.Body).Decode(&student)
	student.ID = strconv.Itoa(rand.Intn(100000000))
	students = append(students, student)
	json.NewEncoder(w).Encode(student)
}
