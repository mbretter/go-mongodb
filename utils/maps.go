package utils

import "go.mongodb.org/mongo-driver/v2/bson"

func Map2BsonM[V any](m map[string]V) bson.M {
	result := make(bson.M, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}
