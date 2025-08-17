// Package types provides the UUID datatype, which derives from string for easy conversion, it's BSON represenation is
// primitive.Binary with the subtype of bson.TypeBinaryUUID.
// This means it is store as native UUID into the database. An empty UUID is treated as null when converting to BSON.
package types

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UUID string

var uuidGenerator = uuid.NewString

// SetUuidGenerator sets a custom function for generating UUID strings.
// This is mainly used for testing purposes.
func SetUuidGenerator(g func() string) {
	uuidGenerator = g
}

// NewUuid generates a new UUID using the configured uuidGenerator function and returns it as a UUID type.
func NewUuid() UUID {
	return UUID(uuidGenerator())
}

// String converts the UUID to its string representation.
func (u UUID) String() string {
	return string(u)
}

// IsZero checks if the UUID is empty, returning true if it is and false otherwise.
func (u UUID) IsZero() bool {
	return len(u) == 0
}

// UuidFromString converts a string representation of a UUID to a UUID type. Returns an error if the string is not a valid UUID.
func UuidFromString(id string) (UUID, error) {
	u, err := uuid.Parse(id)
	return UUID(u.String()), err
}

// MarshalJSON serializes the UUID into a JSON string. If the UUID is empty, it serializes it as null.
func (u UUID) MarshalJSON() ([]byte, error) {
	if u.IsZero() {
		return json.Marshal(nil)
	}

	return json.Marshal(string(u))
}

// UnmarshalJSON deserializes JSON data into the UUID. It handles both non-null and null cases appropriately.
func (u *UUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*u = ""
		return nil
	}

	uid, err := uuid.ParseBytes(data)
	if err != nil {
		return err
	}

	*u = UUID(uid.String())
	return nil
}

// MarshalBSONValue marshals the UUID into a BSON value. Returns BSON type, byte slice, and an error if any.
func (u UUID) MarshalBSONValue() (byte, []byte, error) {
	if u.IsZero() {
		return byte(bson.TypeNull), nil, nil
	}

	uid, err := uuid.Parse(string(u))
	if err != nil {
		return 0, nil, err
	}

	data, err := uid.MarshalBinary()
	if err != nil {
		return 0, nil, err
	}

	bin := bson.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    data,
	}
	return marshalBsonValue(bin)
}

// UnmarshalBSONValue deserializes a BSON value into a UUID. Returns an error if the BSON type or subtype is incorrect.
func (u *UUID) UnmarshalBSONValue(typ byte, data []byte) error {
	t := bson.Type(typ)
	if t == bson.TypeNull {
		*u = ""
		return nil
	}

	if t != bson.TypeBinary {
		return errors.New("wrong bson type expected binary")
	}

	bin := bson.Binary{}

	err := bson.UnmarshalValue(t, data, &bin)
	if err != nil {
		return err
	}

	if bin.Subtype != bson.TypeBinaryUUID {
		return errors.New("wrong subtype")
	}

	uid, err := uuid.FromBytes(bin.Data)
	if err != nil {
		return err
	}

	*u = UUID(uid.String())

	return nil
}
