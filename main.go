package mongoerrorparser

import (
	"encoding/json"
	"errors"

	"github.com/gsdenys/cerr"
	"go.mongodb.org/mongo-driver/mongo"
)

func toMongoException(err error) mongo.WriteException {
	var mongoException mongo.WriteException
	errors.As(err, &mongoException)

	return mongoException
}

func toErrorInterface(myError mongo.WriteException) interface{} {
	var jsonError interface{}

	json.Unmarshal([]byte(myError.Raw.String()), &jsonError)

	return jsonError
}

func Parse(err error) []*cerr.ValidationError {
	ex := toMongoException(err)

	switch ex.WriteConcernError.Code {
	case 121:
		e121 := NewE121(toErrorInterface(ex))
		e121.Run()

		return e121.Errors
	case 11000:
		e11000 := NewE11000(toErrorInterface(ex))
		e11000.Run()

		return e11000.Errors
	}

	return nil
}
