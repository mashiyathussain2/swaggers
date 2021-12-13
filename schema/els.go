package schema

type GetBrandSchema struct {
	ID   string   `json:"_id"`
	Slug string   `json:"slug"`
	Name string   `json:"name"`
	Logo *ImgResp `json:"logo"`
}

type IsLikedByUserBucketSchema struct {
	DocCount uint `json:"doc_count"`
}

type LikeCountBucketSchema struct {
	Key           string                      `json:"key"`
	DocCount      uint                        `json:"doc_count"`
	IsLikedByUser []IsLikedByUserBucketSchema `json:"is_liked_by_user"`
}

type LikeCountSchema struct {
	Buckets []LikeCountBucketSchema `json:"buckets"`
}

type LikeCountAggResp struct {
	LikeCount []LikeCountSchema `json:"like_count"`
}
