// The meparser package provide a way to parse the Mongo response errors. Their
// usage is really simple and consiste just in invocate the function Parse passing
// as parameter the error received by the mongodb operation. For more information
// see the Parse documentation bellow.
package meparser

import (
	"encoding/json"
	"errors"

	"github.com/gsdenys/cerr"
	"go.mongodb.org/mongo-driver/mongo"
)

// validation data structure has a main objective to helps make validation
// validation over the mongo BSON error and provide a model more adherent to the
// golang error interface.
// As the focus of this data structure is be used to parse the validation data,
// the result is an array of ValidationError, that implements the Error() method
// returning a string in JSON format.
type validation struct {
	Errors []error
	Runner []interface{}
}

// toMongoException function convert the error interface in a mongo WriteException
func toMongoException(err error) mongo.WriteException {
	var mongoException mongo.WriteException
	errors.As(err, &mongoException)

	return mongoException
}

// toErrorInterface function convert the mongo WriteException in a JSON object
func toErrorInterface(myError mongo.WriteException) interface{} {
	var jsonError interface{}

	json.Unmarshal([]byte(myError.Raw.String()), &jsonError)

	return jsonError
}

// Execute the mongo error parser. This function receives as parameter the error
// received by the mongodb operation and returns a formatted error array that
// can printed as JSON String.
//
// Note that to facilitate the usage, the status of the returned errors is mapped
// as a HTTP status.
//
// By now, this just parse the error 121 -> DocumentValidationFailure and  the
// 11000 -> DuplicateKey. All other will returns a generic internal server error
// message with status 500.
func Parse(err error) []error {
	ex := toMongoException(err)

	switch ex.WriteConcernError.Code {
	case 121:
		return newE121(toErrorInterface(ex)).Run().Errors
	case 11000:
		return newE11000(toErrorInterface(ex)).Run().Errors
	default:
		return []error{
			cerr.CreateDefaultError("internal server error").Status(500),
		}
	}
}
