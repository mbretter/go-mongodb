// Package types provides the NullString datatype, which encodes empty strings to null and vice versa.
package types

import (
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NullTime time.Time

// MarshalJSON serializes the NullString value to JSON. If the value is empty, it returns JSON null.
func (v NullTime) MarshalJSON() ([]byte, error) {
	if time.Time(v).IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(time.Time(v))
}

// MarshalBSONValue serializes the NullString value to BSON. If the value is empty, it returns BSON null.
func (v NullTime) MarshalBSONValue() (byte, []byte, error) {
	if time.Time(v).IsZero() {
		return byte(bson.TypeNull), nil, nil
	}
	return marshalBsonValue(time.Time(v))
}

// UnmarshalBSONValue deserializes a BSON value into a time.Time.
func (v *NullTime) UnmarshalBSONValue(typ byte, data []byte) error {
	t := bson.Type(typ)
	if t == bson.TypeNull {
		*v = NullTime(time.Time{})
		return nil
	}

	if t != bson.TypeDateTime {
		return errors.New("wrong bson type expected datetime")
	}

	var dt bson.DateTime

	err := bson.UnmarshalValue(t, data, &dt)
	if err != nil {
		return err
	}

	*v = NullTime(dt.Time())

	return nil
}
