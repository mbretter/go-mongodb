package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectId holds a mongodb objectid
// renders to null if zero
type ObjectId string

var NilObjectID = ObjectId("000000000000000000000000")

var objectIdGenerator = func() string {
	newOid := primitive.NewObjectID()
	return newOid.Hex()
}

func SetObjectIdGenerator(g func() string) {
	objectIdGenerator = g
}

func NewObjectId() ObjectId {
	return ObjectId(objectIdGenerator())
}

func ObjectIdFromHex(s string) (ObjectId, error) {
	oId, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return "", err
	}

	return ObjectId(oId.Hex()), nil
}

func (o ObjectId) IsZero() bool {
	return len(o) == 0
}

func (o ObjectId) String() string {
	if o.IsZero() {
		return "null"
	}

	return fmt.Sprintf("ObjectID(%s)", string(o))
}

func (o ObjectId) MarshalJSON() ([]byte, error) {
	if len(o) == 0 {
		return json.Marshal(nil)
	}

	if o == NilObjectID {
		return json.Marshal(nil)
	}

	return json.Marshal(string(o))
}

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

	oId, err := primitive.ObjectIDFromHex(hexStr)
	if err != nil {
		return err
	}

	*o = ObjectId(oId.Hex())
	return nil
}

func (o ObjectId) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if len(o) == 0 {
		return bson.TypeNull, nil, nil
	}
	oId, err := primitive.ObjectIDFromHex(string(o))
	if err != nil {
		return 0, nil, err
	}

	if oId.IsZero() {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(oId)
}

func (o *ObjectId) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t == bson.TypeNull {
		*o = ""
		return nil
	}

	if t != bson.TypeObjectID {
		return errors.New("wrong bson type expected objectid")
	}

	oId := primitive.ObjectID(data)

	*o = ObjectId(oId.Hex())

	return nil
}
