package schema

// Img contains image src url
type Img struct {
	SRC string `json:"src" validate:"required,url"`
}

// ImgResp contains img response info
type ImgResp struct {
	SRC    string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
