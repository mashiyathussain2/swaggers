package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

//CI = Commission Invoice

type GenerateCIEvent struct {
	DebitRequestID primitive.ObjectID `json:"debet_request_id"`
}
