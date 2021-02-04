package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// contains list of db collections
const (
	CategoryColl string = "category"
)

// Category contains single category information such as id name thumbnail etc.
/* Each category can be linked to other category as a children by setting up parent_id which accepts
another category ID. Also Ancestors store the category hierarchy by storing object parent_id of parent until its nil.
*/
type Category struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Slug string             `json:"slug,omitempty" bson:"slug,omitempty"`

	// ParentID stores the direct parent category of the current category
	ParentID primitive.ObjectID `json:"parent_id,omitempty" bson:"parent_id,omitempty"`

	// AncestorID stores the parentID hierarchy; parent of parent until no parent is found
	// IDs of ancestors are stores in reverse order i.e topmost parent (category with no parent first) then descending down to
	// parent id which is the direct parent of the category as the last element.
	AncestorID []primitive.ObjectID `json:"ancestors_id,omitempty" bson:"ancestors_id,omitempty"`

	// thumbnail stores an image which is displayed on the main menu if displayed
	Thumbnail *IMG `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	// FeaturedImage stores a featured image which is displayed on category list view as a heighlight
	FeaturedImage *IMG `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	// IsMain is set to true if you want to display this category in the category-menu
	IsMain bool `json:"is_main,omitempty" bson:"is_main,omitempty"`
}

// Menu contains all the categories and sub-categories that are shown in main menu
type Menu struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	GenderCode []string           `json:"gender_code,omitempty" bson:"gender_code,omitempty"`
	IsActive   bool               `json:"is_active,omitempty" bson:"is_active,omitempty"`
}

// MenuCategory contains category infomation with all its children as well
type MenuCategory struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Slug string             `json:"slug,omitempty" bson:"slug,omitempty"`
	// thumbnail stores an image which is displayed on the main menu if displayed
	Thumbnail *IMG `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	// FeaturedImage stores a featured image which is displayed on category list view as a heighlight
	FeaturedImage *IMG           `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	Children      []MenuCategory `json:"children,omitempty" bson:"children,omitempty"`
}
