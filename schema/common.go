package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

// Img contains image src url

// swagger:parameters Img
type Img struct {
	// required:true
	SRC string `json:"src" validate:"required,url"`
}

type KafkaMeta struct {
	ID        interface{}         `bson:"_id,omitempty" json:"_id,omitempty"`
	Timestamp primitive.Timestamp `bson:"ts" json:"ts"`
	Namespace string              `bson:"ns" json:"ns"`
	Operation string              `bson:"op,omitempty" json:"op,omitempty"`
	Updates   interface{}         `bson:"updates,omitempty" json:"updates,omitempty"`
}

type KafkaMessage struct {
	Meta KafkaMeta              `bson:"meta" json:"meta"`
	Data map[string]interface{} `bson:"data,omitempty" json:"data,omitempty"`
}
