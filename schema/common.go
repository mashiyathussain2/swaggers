package schema

type img struct {
	SRC string `json:"src" validate:"required,url"`
}
