package handlers

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPrimitiveObjectIDFromInterface(in interface{}) (primitive.ObjectID, error) {
	// convert to string and trim the double quotes
	str := strings.Trim(fmt.Sprintf("%s", in), "\"")

	// convert string to primitive Object ID
	oid, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return oid, err
	}

	return oid, nil
}
