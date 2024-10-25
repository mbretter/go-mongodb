package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type NullString string

func (v NullString) MarshalJSON() ([]byte, error) {
	if len(v) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(string(v))
}
func (v NullString) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if len(v) == 0 {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(string(v))
}
