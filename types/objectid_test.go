package types_test

import (
	"encoding/json"
	"github.com/mbretter/go-mongodb/v2/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
)

type ObjectIdTest struct {
	Id types.ObjectId `json:"_id" bson:"_id"`
}

func TestObjectId_New(t *testing.T) {
	oId := types.NewObjectId()
	oId2 := types.NewObjectId()

	emptyOid := types.ObjectId("")

	assert.NotEqual(t, oId, oId2)
	assert.NotZero(t, oId)
	assert.Zero(t, emptyOid)
}

func TestObjectId_NewGenerator(t *testing.T) {
	types.SetObjectIdGenerator(func() string {
		return "6555d2cc4fce49f464c2f683"
	})

	s := types.NewObjectId()
	assert.Equal(t, "6555d2cc4fce49f464c2f683", string(s))
}

func TestObjectId_MarshalJSON(t *testing.T) {
	s := ObjectIdTest{}

	o, _ := types.ObjectIdFromHex("6555d2cc4fce49f464c2f683")
	s.Id = o
	j, _ := json.Marshal(s)

	assert.Equal(t, `{"_id":"6555d2cc4fce49f464c2f683"}`, string(j))
}

func TestObjectId_FromHexInvalid(t *testing.T) {
	oId, err := types.ObjectIdFromHex("x")

	assert.NotNil(t, err)
	assert.Zero(t, oId)
}

func TestObjectId_String(t *testing.T) {
	oId, err := types.ObjectIdFromHex("6555d2cc4fce49f464c2f683")

	assert.Nil(t, err)
	assert.Equal(t, "ObjectID(6555d2cc4fce49f464c2f683)", oId.String())

	oId = types.ObjectId("")
	assert.Equal(t, "null", oId.String())
}

func TestObjectId_MarshalJSONNull(t *testing.T) {
	j, _ := json.Marshal(ObjectIdTest{})

	assert.Equal(t, `{"_id":null}`, string(j))
}

func TestObjectId_MarshalJSONZeroObjectId(t *testing.T) {
	j, _ := json.Marshal(ObjectIdTest{
		Id: "000000000000000000000000",
	})

	assert.Equal(t, `{"_id":null}`, string(j))
}

func TestObjectId_UnmarshalJSON(t *testing.T) {
	s := ObjectIdTest{}

	err := json.Unmarshal([]byte(`{"_id":"6555d2cc4fce49f464c2f683"}`), &s)

	assert.Nil(t, err)
	assert.Equal(t, "6555d2cc4fce49f464c2f683", string(s.Id))
}

func TestObjectId_UnmarshalJSONNull(t *testing.T) {
	s := ObjectIdTest{}

	err := json.Unmarshal([]byte(`{"_id":null}`), &s)

	assert.Nil(t, err)
	assert.Equal(t, "", string(s.Id))
}

func TestObjectId_UnmarshalJSONInvalidObjectId(t *testing.T) {
	s := ObjectIdTest{}

	err := json.Unmarshal([]byte(`{"_id":"f47ac10b"}`), &s)
	if assert.NotNil(t, err) {
		assert.Equal(t, "the provided hex string is not a valid ObjectID", err.Error())
	}
}

func TestObjectId_UnmarshalJSONInvalidJson(t *testing.T) {
	s := ObjectIdTest{}

	err := s.Id.UnmarshalJSON([]byte(`"`))
	if assert.NotNil(t, err) {
		assert.Equal(t, "unexpected end of JSON input", err.Error())
	}
}

func TestObjectId_MarshalBSON(t *testing.T) {
	s := ObjectIdTest{}

	u, _ := types.ObjectIdFromHex("6555d2cc4fce49f464c2f683")
	s.Id = u
	b, err := bson.Marshal(s)

	assert.Nil(t, err)
	assert.Equal(t, "\x16\x00\x00\x00\a_id\x00eU\xd2\xccO\xceI\xf4d\xc2\xf6\x83\x00", string(b))
}

func TestObjectId_MarshalBSONInvalid(t *testing.T) {
	s := ObjectIdTest{
		Id: "xxx",
	}

	_, err := bson.Marshal(s)

	assert.NotNil(t, err)
	assert.Equal(t, "the provided hex string is not a valid ObjectID", err.Error())
}

func TestObjectId_MarshalBSONNull(t *testing.T) {
	s := ObjectIdTest{}

	b, err := bson.Marshal(s)

	assert.Nil(t, err)
	assert.Equal(t, "\n\x00\x00\x00\n_id\x00\x00", string(b))
}

func TestObjectId_MarshalBSONZeroObjectId(t *testing.T) {
	s := ObjectIdTest{}

	u, _ := types.ObjectIdFromHex("000000000000000000000000")
	s.Id = u
	b, err := bson.Marshal(s)

	assert.Nil(t, err)
	assert.Equal(t, "\n\x00\x00\x00\n_id\x00\x00", string(b))
}

func TestObjectId_UnmarshalBSON(t *testing.T) {
	s := ObjectIdTest{}

	err := bson.Unmarshal([]byte("\x16\x00\x00\x00\a_id\x00eU\xd2\xccO\xceI\xf4d\xc2\xf6\x83\x00"), &s)

	assert.Nil(t, err)
	assert.Equal(t, "6555d2cc4fce49f464c2f683", string(s.Id))
}

func TestObjectId_UnmarshalBSONNull(t *testing.T) {
	s := ObjectIdTest{}

	err := bson.Unmarshal([]byte("\n\x00\x00\x00\n_id\x00\x00"), &s)

	assert.Nil(t, err)
	assert.Equal(t, "", string(s.Id))
}

func TestObjectId_UnmarshalBSONWrongType(t *testing.T) {
	s := ObjectIdTest{}

	err := bson.Unmarshal([]byte("\x1f\x00\x00\x00\x05_id\x00\x10\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), &s)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding key _id: wrong bson type expected objectid", err.Error())
}
