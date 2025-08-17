package types

import "go.mongodb.org/mongo-driver/v2/bson"

func marshalBsonValue(data any) (byte, []byte, error) {
	typ, v, err := bson.MarshalValue(data)
	return byte(typ), v, err
}
