// Package types provides an ObjectId replacement of the mongodb drivers primitive.ObjectId.
// The original ObjectId has two disadvantages:
//
// * an empty ObjectId is stored as "000000000000000000000000", instead of null, which is kind of weird.
//
// * every conversion between a string and an ObjectId has to be done using ObjectIDFromHex(), which adds a lot of extra code.
package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ObjectId holds a mongodb objectid
// renders to null if zero
type ObjectId string

var NilObjectID = ObjectId("000000000000000000000000")

var objectIdGenerator = func() string {
	newOid := bson.NewObjectID()
	return newOid.Hex()
}

// SetObjectIdGenerator sets a custom function for generating ObjectId strings.
// This is mainly used for testing purposes.
func SetObjectIdGenerator(g func() string) {
	objectIdGenerator = g
}

// NewObjectId generates and returns a new ObjectId using the objectIdGenerator function.
func NewObjectId() ObjectId {
	return ObjectId(objectIdGenerator())
}

// ObjectIdFromHex converts a hexadecimal string representation of an ObjectId to an ObjectId type, returning an error if invalid.
func ObjectIdFromHex(s string) (ObjectId, error) {
	oId, err := bson.ObjectIDFromHex(s)
	if err != nil {
		return "", err
	}

	return ObjectId(oId.Hex()), nil
}

// IsZero checks if the ObjectId is zero (an empty string) and returns true if it is, otherwise false.
func (o ObjectId) IsZero() bool {
	return len(o) == 0
}

// String returns the string representation of the ObjectId. If the ObjectId is zero, it returns "null".
func (o ObjectId) String() string {
	if o.IsZero() {
		return "null"
	}

	return fmt.Sprintf("ObjectID(%s)", string(o))
}

// MarshalJSON serializes the ObjectId to JSON, rendering as null if the ObjectId is zero or equals NilObjectID.
func (o ObjectId) MarshalJSON() ([]byte, error) {
	if len(o) == 0 {
		return json.Marshal(nil)
	}

	if o == NilObjectID {
		return json.Marshal(nil)
	}

	return json.Marshal(string(o))
}

// UnmarshalJSON unmarshals a JSON-encoded string into an ObjectId. Handles null values by setting the ObjectId to an empty string.
func (o *ObjectId) UnmarshalJSON(data []byte) error {
	var hexStr string
	err := json.Unmarshal(data, &hexStr)
	if err != nil {
		return err
	}

	// null value
	if len(hexStr) == 0 {
		*o = ""
		return nil
	}

	oId, err := bson.ObjectIDFromHex(hexStr)
	if err != nil {
		return err
	}

	*o = ObjectId(oId.Hex())
	return nil
}

// MarshalBSONValue serializes the ObjectId to BSON. It returns BSON null type for zero values or in case of invalid ObjectId.
func (o ObjectId) MarshalBSONValue() (byte, []byte, error) {
	if len(o) == 0 {
		return byte(bson.TypeNull), nil, nil
	}
	oId, err := bson.ObjectIDFromHex(string(o))
	if err != nil {
		return 0, nil, err
	}

	if oId.IsZero() {
		return byte(bson.TypeNull), nil, nil
	}

	return marshalBsonValue(oId)
}

// UnmarshalBSONValue deserializes BSON value into ObjectId type. Returns an error if the BSON type is invalid.
func (o *ObjectId) UnmarshalBSONValue(typ byte, data []byte) error {
	t := bson.Type(typ)
	if t == bson.TypeNull {
		*o = ""
		return nil
	}

	if t != bson.TypeObjectID {
		return errors.New("wrong bson type expected objectid")
	}

	oId := bson.ObjectID(data)

	*o = ObjectId(oId.Hex())

	return nil
}
