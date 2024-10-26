// Package types provides various number datatypes, they are treated as BSON-null if their value is 0 oder 0.0 and vice versa.
package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type NullFloat32 float64

// MarshalJSON serializes the NullFloat32 value to JSON. If the value is 0, it marshals to JSON null. Otherwise, it marshals as float32.
func (v NullFloat32) MarshalJSON() ([]byte, error) {
	if v == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(float32(v))
}

// MarshalBSONValue serializes the NullFloat32 value to BSON. If the value is 0, it marshals to BSON null; otherwise, as float32.
func (v NullFloat32) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if v == 0 {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(float32(v))
}

type NullFloat64 float64

// MarshalJSON customizes the JSON marshaling process for NullFloat64. It returns nil if the value is 0, otherwise it returns the float64 value.
func (v NullFloat64) MarshalJSON() ([]byte, error) {
	if v == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(float64(v))
}

// MarshalBSONValue customizes the BSON marshaling process for NullFloat64. It returns BSON Null type
func (v NullFloat64) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if v == 0 {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(float64(v))
}

type NullInt32 int32

// MarshalJSON serializes the NullInt32 value into JSON, encoding zero values as null.
func (v NullInt32) MarshalJSON() ([]byte, error) {
	if v == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(int32(v))
}

// MarshalBSONValue serializes the NullInt32 value into BSON format, encoding zero values as null.
func (v NullInt32) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if v == 0 {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(int32(v))
}

type NullInt64 int64

// MarshalJSON marshals the NullInt64 value into JSON. If the value is zero, it marshals as null.
func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(int64(v))
}

// MarshalBSONValue serializes NullInt64 into a BSON value. Returns BSON null type if the value is zero.
func (v NullInt64) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if v == 0 {
		return bson.TypeNull, nil, nil
	}
	return bson.MarshalValue(int64(v))
}
