package meparser

import (
	"fmt"

	"github.com/gsdenys/cerr"
)

// e11000 data structure to parser the mongodb validation error for error
// code 121
type e11000 struct {
	validation
}

// newE11000 function that provide a way to create an E11000 instance
func newE11000(root interface{}) *e11000 {
	e := &e11000{}
	e.Runner = append(e.Runner, root)

	return e
}

// getWriteErrors function to return a map structured containing the field name
// and their error
func (e *e11000) getWriteErrors() map[string]interface{} {
	mapa := e.Runner[0].(map[string]interface{})
	writeErrors := mapa["writeErrors"].([]interface{})[0]

	return writeErrors.(map[string]interface{})
}

// Run function to execute the parser over the E11000 mongodb error. This function
// perform an parser over BSON error and store all errors inside their Errors
// data structure.
func (e *e11000) Run() *e11000 {
	we := e.getWriteErrors()

	writeErros := we["keyValue"].(map[string]interface{})

	for k, v := range writeErros {
		message := fmt.Sprintf("%s already exist and cannot be duplicated", v.(string))
		e.Errors = append(
			e.Errors,
			cerr.CreateValidationError(k, message).Status(422),
		)
	}

	return e
}
