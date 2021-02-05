package schema

// Img contains image src url
type Img struct {
	SRC string `json:"src" validate:"required,url"`
}
