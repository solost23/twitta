package utils

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParseObjectID converts ObjectID("hexstring") or plain hex to primitive.ObjectID.
func ParseObjectID(s string) (primitive.ObjectID, error) {
	s = strings.TrimPrefix(s, `ObjectID("`)
	s = strings.TrimSuffix(s, `")`)
	return primitive.ObjectIDFromHex(s)
}
