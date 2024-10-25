package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type StringTest struct {
	Name NullString `json:"name"`
}

func TestString_MarshalEmpty(t *testing.T) {
	s := StringTest{}

	j, _ := json.Marshal(s)
	b, _ := bson.Marshal(s)

	assert.Equal(t, `{"name":null}`, string(j))
	assert.Equal(t, "\v\x00\x00\x00\nname\x00\x00", string(b))
}

func TestString_MarshalNonEmpty(t *testing.T) {
	s := StringTest{
		Name: "foo",
	}

	j, _ := json.Marshal(s)
	b, _ := bson.Marshal(s)

	assert.Equal(t, `{"name":"foo"}`, string(j))
	assert.Equal(t, "\x13\x00\x00\x00\x02name\x00\x04\x00\x00\x00foo\x00\x00", string(b))

}
