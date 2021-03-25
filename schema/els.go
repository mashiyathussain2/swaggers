package schema

type GetBrandSchema struct {
	ID   string   `json:"_id"`
	Slug string   `json:"slug"`
	Name string   `json:"name"`
	Logo *ImgResp `json:"logo"`
}
