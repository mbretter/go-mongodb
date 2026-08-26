package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mbretter/go-mongodb/v2/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TimeTest struct {
	Created types.NullTime `json:"created" bson:"created"`
}

func TestTime_MarshalJSON(t *testing.T) {
	s := TimeTest{}

	u, _ := time.Parse(time.RFC3339, "2021-01-01T10:02:30Z")
	s.Created = types.NullTime(u)
	j, _ := json.Marshal(s)

	assert.Equal(t, `{"created":"2021-01-01T10:02:30Z"}`, string(j))
}

func TestTime_MarshalBSON(t *testing.T) {
	s := TimeTest{}

	u, _ := time.Parse(time.RFC3339, "2021-01-01T10:02:30Z")
	s.Created = types.NullTime(u)
	b, err := bson.Marshal(s)

	assert.Nil(t, err)
	assert.Equal(t, "\x16\x00\x00\x00\tcreated\x00\xf0\nf\xbdv\x01\x00\x00\x00", string(b))
}

func TestTime_MarshalBSONNull(t *testing.T) {
	s := TimeTest{}

	b, err := bson.Marshal(s)

	assert.Nil(t, err)
	assert.Equal(t, "\x0e\x00\x00\x00\ncreated\x00\x00", string(b))
}

func TestTime_UnmarshalBSONNull(t *testing.T) {
	s := TimeTest{}

	err := bson.Unmarshal([]byte("\x0e\x00\x00\x00\ncreated\x00\x00"), &s)

	assert.Nil(t, err)
	assert.Zero(t, s.Created)
}
