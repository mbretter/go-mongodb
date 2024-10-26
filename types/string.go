// Package types provides the NullString datatype, which encodes empty strings to null and vice versa.
package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type NullString string

// MarshalJSON serializes the NullString value to JSON. If the value is empty, it returns JSON null.
func (v NullString) MarshalJSON() ([]byte, error) {
	if len(v) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(string(v))
}

// MarshalBSONValue serializes the NullString value to BSON. If the value is empty, it returns BSON null.
func (v NullString) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if len(v) == 0 {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(string(v))
}
