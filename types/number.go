package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type NullFloat32 float64

func (v NullFloat32) MarshalJSON() ([]byte, error) {
	if v == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(float32(v))
}
func (v NullFloat32) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if v == 0 {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(float32(v))
}

type NullFloat64 float64

func (v NullFloat64) MarshalJSON() ([]byte, error) {
	if v == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(float64(v))
}
func (v NullFloat64) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if v == 0 {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(float64(v))
}

type NullInt32 int32

func (v NullInt32) MarshalJSON() ([]byte, error) {
	if v == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(int32(v))
}
func (v NullInt32) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if v == 0 {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(int32(v))
}

type NullInt64 int64

func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(int64(v))
}
func (v NullInt64) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if v == 0 {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(int64(v))
}
