package mongoerrorparser

import (
	"fmt"

	ce "github.com/finacore/commons-errors"
)

type E11000 struct {
	Errors []*ce.ValidationError
	Runner []interface{}
}

func NewE11000(root interface{}) *E11000 {
	e := &E11000{}
	e.Runner = append(e.Runner, root)

	return e
}

func (e *E11000) getWriteErrors() map[string]interface{} {
	mapa := e.Runner[0].(map[string]interface{})
	writeErrors := mapa["writeErrors"].([]interface{})[0]

	return writeErrors.(map[string]interface{})
}

func (e *E11000) Run() {
	we := e.getWriteErrors()

	writeErros := we["keyValue"].(map[string]interface{})

	for k, v := range writeErros {
		message := fmt.Sprintf("%s already exist and cannot be duplicated", v.(string))
		e.Errors = append(e.Errors, ce.CreateValidationError(k, message))
	}
}
