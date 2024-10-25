package types

import (
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type BinaryTest struct {
	Data Binary `json:"data" bson:"data"`
}

// png 10x10 white bg
const binaryTestImage = `iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKAQAAAAClSfIQAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QAAd2KE6QAAAAHdElNRQfnDAEKNi+VgyIEAAAADklEQVQI12P4f4ABNwIAB1IRd+bI0OMAAAAldEVYdGRhdGU6Y3JlYXRlADIwMjMtMTItMDFUMTA6NTQ6NDcrMDA6MDAuAbquAAAAJXRFWHRkYXRlOm1vZGlmeQAyMDIzLTEyLTAxVDEwOjU0OjQ3KzAwOjAwX1wCEgAAAABJRU5ErkJggg==`

var binaryTestData, _ = base64.StdEncoding.DecodeString(binaryTestImage)

func TestBinary_New(t *testing.T) {
	s := BinaryTest{}

	assert.Nil(t, s.Data)
	assert.Len(t, s.Data, 0)

	s.Data = []byte{}
	assert.NotNil(t, s.Data)
	assert.Len(t, s.Data, 0)
}

func TestBinary_MarshalJSONEmpty(t *testing.T) {
	s := BinaryTest{}

	j, _ := json.Marshal(s)
	assert.Equal(t, `{"data":null}`, string(j))

	s.Data = []byte{}
	j, _ = json.Marshal(s)
	assert.Equal(t, `{"data":null}`, string(j))
}

func TestBinary_MarshalJSON(t *testing.T) {
	s := BinaryTest{
		Data: binaryTestData,
	}

	j, _ := json.Marshal(s)
	assert.Equal(t, `{"data":"`+binaryTestImage+`"}`, string(j))
}

func TestBinary_UnmarshalJSON(t *testing.T) {
	s := BinaryTest{}

	jsonStr := `{"data":"` + binaryTestImage + `"}`

	_ = json.Unmarshal([]byte(jsonStr), &s)
	assert.Equal(t, binaryTestData, []byte(s.Data))
}

func TestBinary_UnmarshalJSONEmpty(t *testing.T) {
	s := BinaryTest{}

	_ = json.Unmarshal([]byte(`{"data":""}`), &s)
	assert.Equal(t, Binary(nil), s.Data)

	_ = json.Unmarshal([]byte(`{"data":null}`), &s)
	assert.Equal(t, Binary(nil), s.Data)

	err := s.Data.UnmarshalJSON([]byte(""))
	if assert.Nil(t, err) {
		assert.Nil(t, s.Data, Binary(nil))
	}
}

func TestBinary_UnmarshalJSONInvalidBase64(t *testing.T) {
	s := BinaryTest{}

	err := s.Data.UnmarshalJSON([]byte("x"))
	if assert.NotNil(t, err) {
		assert.Equal(t, "invalid character 'x' looking for beginning of value", err.Error())
	}
}

func TestBinary_MarshalBSON(t *testing.T) {
	s := BinaryTest{
		Data: binaryTestData,
	}

	b, err := bson.Marshal(s)

	if assert.Nil(t, err) {
		assert.Equal(t, "\x16\x01\x00\x00\x05data\x00\x06\x01\x00\x00\x00\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\n\x00\x00\x00\n\x01\x00\x00\x00\x00\xa5I\xf2\x10\x00\x00\x00\x04gAMA\x00\x00\xb1\x8f\v\xfca\x05\x00\x00\x00 cHRM\x00\x00z&\x00\x00\x80\x84\x00\x00\xfa\x00\x00\x00\x80\xe8\x00\x00u0\x00\x00\xea`\x00\x00:\x98\x00\x00\x17p\x9c\xbaQ<\x00\x00\x00\x02bKGD\x00\x01݊\x13\xa4\x00\x00\x00\atIME\a\xe7\f\x01\n6/\x95\x83\"\x04\x00\x00\x00\x0eIDAT\b\xd7c\xf8\x7f\x80\x017\x02\x00\aR\x11w\xe6\xc8\xd0\xe3\x00\x00\x00%tEXtdate:create\x002023-12-01T10:54:47+00:00.\x01\xba\xae\x00\x00\x00%tEXtdate:modify\x002023-12-01T10:54:47+00:00_\\\x02\x12\x00\x00\x00\x00IEND\xaeB`\x82\x00", string(b))
	}
}

func TestBinary_MarshalBSONEmpty(t *testing.T) {
	s := BinaryTest{}

	b, err := bson.Marshal(s)

	if assert.Nil(t, err) {
		assert.Equal(t, "\v\x00\x00\x00\ndata\x00\x00", string(b))
	}
}

func TestBinary_UnmarshalBSON(t *testing.T) {
	s := BinaryTest{}

	err := bson.Unmarshal([]byte("\x16\x01\x00\x00\x05data\x00\x06\x01\x00\x00\x00\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\n\x00\x00\x00\n\x01\x00\x00\x00\x00\xa5I\xf2\x10\x00\x00\x00\x04gAMA\x00\x00\xb1\x8f\v\xfca\x05\x00\x00\x00 cHRM\x00\x00z&\x00\x00\x80\x84\x00\x00\xfa\x00\x00\x00\x80\xe8\x00\x00u0\x00\x00\xea`\x00\x00:\x98\x00\x00\x17p\x9c\xbaQ<\x00\x00\x00\x02bKGD\x00\x01݊\x13\xa4\x00\x00\x00\atIME\a\xe7\f\x01\n6/\x95\x83\"\x04\x00\x00\x00\x0eIDAT\b\xd7c\xf8\x7f\x80\x017\x02\x00\aR\x11w\xe6\xc8\xd0\xe3\x00\x00\x00%tEXtdate:create\x002023-12-01T10:54:47+00:00.\x01\xba\xae\x00\x00\x00%tEXtdate:modify\x002023-12-01T10:54:47+00:00_\\\x02\x12\x00\x00\x00\x00IEND\xaeB`\x82\x00"), &s)

	if assert.Nil(t, err) {
		assert.Equal(t, Binary(binaryTestData), s.Data)
	}
}

func TestBinary_UnmarshalBSONEmpty(t *testing.T) {
	s := BinaryTest{}

	err := bson.Unmarshal([]byte("\v\x00\x00\x00\ndata\x00\x00"), &s)

	if assert.Nil(t, err) {
		assert.Equal(t, Binary(nil), s.Data)
	}
}

func TestBinary_UnmarshalBSONInvalidType(t *testing.T) {
	s := BinaryTest{}

	err := bson.Unmarshal([]byte("\x13\x00\x00\x00\x02data\x00\x04\x00\x00\x00foo\x00\x00"), &s)

	if assert.NotNil(t, err) {
		assert.Equal(t, "error decoding key data: wrong bson type expected binary", err.Error())
	}
}

func TestBinary_UnmarshalBSONInvalidSubType(t *testing.T) {
	s := BinaryTest{}

	err := bson.Unmarshal([]byte(" \x00\x00\x00\x05data\x00\x10\x00\x00\x00\x04\x87\x89r^\xb4\aA\x1f\x82`}\x93\xc0\x12k\xc8\x00"), &s)

	if assert.NotNil(t, err) {
		assert.Equal(t, "error decoding key data: wrong bson subtype expected generic", err.Error())
	}
}
