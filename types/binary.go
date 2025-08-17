// Package types provide a binary datatype. The binary datatype stores any arbitrary value as binary,
// the binary subtype is `bson.TypeBinaryGeneric`. The JSON representation of the binary is base64.
package types

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Binary []byte

// MarshalJSON serializes the Binary receiver as a base64-encoded string or null if empty.
func (b Binary) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return json.Marshal(nil)
	}

	return json.Marshal(base64.StdEncoding.EncodeToString(b))
}

// UnmarshalJSON decodes a JSON-encoded byte slice as a base64-encoded string and stores the result in the Binary receiver.
func (b *Binary) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*b = nil
		return nil
	}

	var base64Str string
	err := json.Unmarshal(data, &base64Str)
	if err != nil {
		return err
	}

	if len(base64Str) == 0 {
		*b = nil
		return nil
	}

	bytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return err
	}

	*b = bytes

	return nil
}

// MarshalBSONValue serializes the Binary receiver into a BSON value.
func (b Binary) MarshalBSONValue() (byte, []byte, error) {
	if len(b) == 0 {
		return byte(bson.TypeNull), nil, nil
	}

	return marshalBsonValue(bson.Binary{Data: b, Subtype: bson.TypeBinaryGeneric})
}

// UnmarshalBSONValue decodes a BSON value into the Binary receiver, ensuring it is of the correct BSON type and subtype.
func (b *Binary) UnmarshalBSONValue(typ byte, data []byte) error {
	t := bson.Type(typ)
	if t == bson.TypeNull {
		*b = nil
		return nil
	}

	if t != bson.TypeBinary {
		return errors.New("wrong bson type expected binary")
	}

	prim := bson.Binary{}
	err := bson.UnmarshalValue(t, data, &prim)
	if err != nil {
		return err
	}

	if prim.Subtype != bson.TypeBinaryGeneric {
		return errors.New("wrong bson subtype expected generic")
	}

	*b = prim.Data

	return nil
}
