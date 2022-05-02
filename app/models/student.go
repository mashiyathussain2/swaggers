package models

// swagger:parameters error commonError
type CommonError struct {
	// Status of the error
	Status int `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:model student createStudent
type Student struct {
	// Id of the student
	// in: string
	ID string `json:"id"`
	// Name of the student
	// in: string
	Name string `json:"name"`
	// Subject of the student
	// in: string
	// required: true
	Subject string `json:"subject"`
}

// swagger:parameters student createStudent
type AddStudentBody struct {
	// - name: body
	//  in: body
	//  description: id, name and subject
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/Student"
	//  required: true
	Body Student `json:"body"`
}

// swagger:parameters error commonError
type AddErrorBody struct {
	// - name: body
	//  in: body
	//  description: name and subject
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/CommonError"
	//  required: true
	Body CommonError `json:"body"`
}
