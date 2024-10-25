package types

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UUID string

var uuidGenerator = uuid.NewString

func SetUuidGenerator(g func() string) {
	uuidGenerator = g
}

func NewUuid() UUID {
	return UUID(uuidGenerator())
}

func (u UUID) String() string {
	return string(u)
}

func (u UUID) IsZero() bool {
	return len(u) == 0
}

func UuidFromString(id string) (UUID, error) {
	u, err := uuid.Parse(id)
	return UUID(u.String()), err
}

func (u UUID) MarshalJSON() ([]byte, error) {
	if u.IsZero() {
		return json.Marshal(nil)
	}

	return json.Marshal(string(u))
}

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

func (u UUID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if u.IsZero() {
		return bson.TypeNull, nil, nil
	}

	uid, err := uuid.Parse(string(u))
	if err != nil {
		return 0, nil, err
	}

	data, err := uid.MarshalBinary()
	if err != nil {
		return 0, nil, err
	}

	bin := primitive.Binary{
		Subtype: bson.TypeBinaryUUID,
		Data:    data,
	}
	return bson.MarshalValue(bin)
}

func (u *UUID) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t == bson.TypeNull {
		*u = ""
		return nil
	}

	if t != bson.TypeBinary {
		return errors.New("wrong bson type expected binary")
	}

	bin := primitive.Binary{}

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
