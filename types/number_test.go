package types_test

import (
	"encoding/json"
	"github.com/mbretter/go-mongodb/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type NumberTest struct {
	Float32 types.NullFloat32 `json:"float32"`
	Float64 types.NullFloat64 `json:"float64"`
	Int32   types.NullInt32   `json:"int32"`
	Int64   types.NullInt64   `json:"int64"`
}

func TestNumber_MarshalEmpty(t *testing.T) {
	s := NumberTest{}

	j, _ := json.Marshal(s)
	b, _ := bson.Marshal(s)

	assert.Equal(t, `{"float32":null,"float64":null,"int32":null,"int64":null}`, string(j))
	assert.Equal(t, "%\x00\x00\x00\nfloat32\x00\nfloat64\x00\nint32\x00\nint64\x00\x00", string(b))

}

func TestNumber_MarshalNonEmpty(t *testing.T) {
	s := NumberTest{
		Int32:   0x7FFFFFFF,
		Int64:   0x7FFFFFFFFFFFFFFF,
		Float32: 1.3,
		Float64: 0.3,
	}

	j, _ := json.Marshal(s)
	b, _ := bson.Marshal(s)

	assert.Equal(t, `{"float32":1.3,"float64":0.3,"int32":2147483647,"int64":9223372036854775807}`, string(j))
	assert.Equal(t, "A\x00\x00\x00\x01float32\x00\x00\x00\x00\xc0\xcc\xcc\xf4?\x01float64\x00333333\xd3?\x10int32\x00\xff\xff\xff\x7f\x12int64\x00\xff\xff\xff\xff\xff\xff\xff\x7f\x00", string(b))

}
