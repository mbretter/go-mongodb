package types_test

import (
	"encoding/json"
	"github.com/mbretter/go-mongodb/v2/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
)

type UuidTest struct {
	Uid types.UUID `json:"uuid" bson:"uuid"`
}

func TestUUID_New(t *testing.T) {
	s := types.NewUuid()

	assert.NotZero(t, s.String())

	s2 := types.UUID("")
	assert.True(t, s2.IsZero())
	assert.Zero(t, s2.String())
}

func TestUUID_NewGenerator(t *testing.T) {
	types.SetUuidGenerator(func() string {
		return "f47ac10b-58cc-0372-8567-0e02b2c3d479"
	})

	s := types.NewUuid()
	assert.Equal(t, "f47ac10b-58cc-0372-8567-0e02b2c3d479", s.String())
}

func TestUUID_MarshalJSON(t *testing.T) {
	s := UuidTest{}

	u, _ := types.UuidFromString("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	s.Uid = u
	j, _ := json.Marshal(s)

	assert.Equal(t, `{"uuid":"f47ac10b-58cc-0372-8567-0e02b2c3d479"}`, string(j))
}

func TestUUID_MarshalJSONNull(t *testing.T) {
	s := UuidTest{}

	j, _ := json.Marshal(s)

	assert.Equal(t, `{"uuid":null}`, string(j))
}

func TestUUID_UnmarshalJSON(t *testing.T) {
	s := UuidTest{}

	err := json.Unmarshal([]byte(`{"uuid":"f47ac10b-58cc-0372-8567-0e02b2c3d479"}`), &s)

	assert.Nil(t, err)
	assert.Equal(t, "f47ac10b-58cc-0372-8567-0e02b2c3d479", s.Uid.String())
}

func TestUUID_UnmarshalJSONNull(t *testing.T) {
	s := UuidTest{}

	err := json.Unmarshal([]byte(`{"uuid":null}`), &s)

	assert.Nil(t, err)
	assert.Equal(t, "", s.Uid.String())
}

func TestUUID_UnmarshalJSONNotExistend(t *testing.T) {
	s := UuidTest{}

	err := json.Unmarshal([]byte(`{}`), &s)

	assert.Nil(t, err)
	assert.Equal(t, "", s.Uid.String())
}

func TestUUID_UnmarshalJSONInvalidUuid(t *testing.T) {
	s := UuidTest{}

	err := json.Unmarshal([]byte(`{"uuid":"f47ac10b-58cc-0372-8567-0e02b2c3d4"}`), &s)

	assert.NotNil(t, err)
	assert.Equal(t, "invalid UUID format", err.Error())
}

func TestUUID_MarshalBSON(t *testing.T) {
	s := UuidTest{}

	u, _ := types.UuidFromString("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	s.Uid = u
	b, err := bson.Marshal(s)

	assert.Nil(t, err)
	assert.Equal(t, " \x00\x00\x00\x05uuid\x00\x10\x00\x00\x00\x04\xf4z\xc1\vX\xcc\x03r\x85g\x0e\x02\xb2\xc3\xd4y\x00", string(b))
}

func TestUUID_MarshalBSONInvalid(t *testing.T) {
	s := UuidTest{
		Uid: "xxxx",
	}
	_, err := bson.Marshal(s)

	assert.NotNil(t, err)
	assert.Equal(t, "invalid UUID length: 4", err.Error())
}

func TestUUID_MarshalBSONNull(t *testing.T) {
	s := UuidTest{}

	b, err := bson.Marshal(s)

	assert.Nil(t, err)
	assert.Equal(t, "\v\x00\x00\x00\nuuid\x00\x00", string(b))
}

func TestUUID_UnmarshalBSON(t *testing.T) {
	s := UuidTest{}

	err := bson.Unmarshal([]byte(" \x00\x00\x00\x05uuid\x00\x10\x00\x00\x00\x04\xf4z\xc1\vX\xcc\x03r\x85g\x0e\x02\xb2\xc3\xd4y\x00"), &s)

	assert.Nil(t, err)
	assert.Equal(t, "f47ac10b-58cc-0372-8567-0e02b2c3d479", s.Uid.String())
}

func TestUUID_UnmarshalBSONNull(t *testing.T) {
	s := UuidTest{}

	err := bson.Unmarshal([]byte("\v\x00\x00\x00\nuuid\x00\x00"), &s)

	assert.Nil(t, err)
	assert.Equal(t, "", s.Uid.String())
}

func TestUUID_UnmarshalBSONWrongType(t *testing.T) {
	s := UuidTest{}

	err := bson.Unmarshal([]byte(" \x00\x00\x00\x04uuid\x00\x10\x00\x00\x00\x03\xf4z\xc1\vX\xcc\x03r\x85g\x0e\x02\xb2\xc3\xd4y\x00"), &s)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding key uuid: wrong bson type expected binary", err.Error())
}

func TestUUID_UnmarshalBSONWrongSubtype(t *testing.T) {
	s := UuidTest{}

	err := bson.Unmarshal([]byte(" \x00\x00\x00\x05uuid\x00\x10\x00\x00\x00\x03\xf4z\xc1\vX\xcc\x03r\x85g\x0e\x02\xb2\xc3\xd4y\x00"), &s)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding key uuid: wrong subtype", err.Error())
}

func TestUUID_UnmarshalBSONInvalidUuid(t *testing.T) {
	s := UuidTest{}

	err := bson.Unmarshal([]byte("\x12\x00\x00\x00\x05uuid\x00\x02\x00\x00\x00\x04\x0e\x02\x00"), &s)

	assert.NotNil(t, err)
	assert.Equal(t, "error decoding key uuid: invalid UUID (got 2 bytes)", err.Error())
}
