package types

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Binary []byte

func (b Binary) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return json.Marshal(nil)
	}

	return json.Marshal(base64.StdEncoding.EncodeToString(b))
}

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

func (b Binary) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if len(b) == 0 {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(primitive.Binary{Data: b, Subtype: bson.TypeBinaryGeneric})
}

func (b *Binary) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t == bson.TypeNull {
		*b = nil
		return nil
	}

	if t != bson.TypeBinary {
		return errors.New("wrong bson type expected binary")
	}

	prim := primitive.Binary{}
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
