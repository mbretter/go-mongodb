package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestMap2Bson(t *testing.T) {
	b1 := Map2BsonM(map[string]string{"a": "b"})
	assert.Equal(t, b1, bson.M{"a": "b"})

	b2 := Map2BsonM(map[string]int{"a": 222})
	assert.Equal(t, b2, bson.M{"a": 222})
}
