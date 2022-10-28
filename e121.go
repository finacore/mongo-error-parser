package meparser

import (
	"github.com/gsdenys/cerr"
)

// e121 data structure to parser the mongodb validation error for error code 121
type e121 struct {
	validation
}

// newE121 allows to create a new E121 data structure prepared to be executed. As
// input parameter, this function receives the interface to be parsed.
func newE121(root interface{}) *e121 {
	e := &e121{}
	e.Runner = append(e.Runner, root)

	return e
}

// deleteElement removes the first element of execution array, providing an
// implementation of a pointer at the breadth search algorithm.
func (e *e121) deleteElement() {
	e.Runner = e.Runner[1:]
}

// addAll enables to add all elements of the pointered element to the end of
// array. Note that this method just be used case the pointered element is an
// implementation off []interface{}.
func (e *e121) addAll() {
	children := e.Runner[0].([]interface{})
	for index := range children {
		e.Runner = append(e.Runner, children[index])
	}
}

// addChildren adds all object's elements to the processing array. Different of
// the addAll, this method requires that the interface be of a
// type map[string]interface{}
func (e *e121) addChildren() {
	data := e.Runner[0].(map[string]interface{})

	for _, v := range data {
		e.Runner = append(e.Runner, v)
	}
}

// addError add a new error to the error array
func (e *e121) addError(field string, message string) {
	e.Errors = append(
		e.Errors,
		cerr.CreateValidationError(field, message).Status(422),
	)
}

// getReason returns the reason stored inside the details map. By now, this
// function just accept one reason for each field.
func (e *e121) getReason() []string {
	var reason []string
	data := e.Runner[0].(map[string]interface{})
	details := data["details"].([]interface{})

	for index := range details {
		switch details[index].(type) {
		case map[string]interface{}:
			item := details[index].(map[string]interface{})
			reason = append(reason, item["reason"].(string))
		}
	}

	return reason
}

// getErrorMessage returns the the error message based on their type. the
// priority is the message that the user has set to the mongodb in the moment of
// create the constraint. In case of this not exist, the reason finding process
// will be started
func (e *e121) getErrorMessage(data map[string]interface{}) []string {
	description := data["description"]

	if description != nil {
		var desc []string
		return append(desc, description.(string))
	}

	return e.getReason()
}

// processInterface verify if the interface is an error container or if it is
// just one more node.
// if this is just one more error, their children ill be add to the process list,
// other else the error will be created and add to the error array
func (e *e121) processInterface() {
	data := e.Runner[0].(map[string]interface{})
	propertyName := data["propertyName"]

	if propertyName == nil {
		e.addChildren()
	} else {
		msgs := e.getErrorMessage(data)

		for i := range msgs {
			e.addError(propertyName.(string), msgs[i])
		}
	}
}

// Run function to execute the parser over the E121 mongodb error. This function
// perform an parser over BSON error and store all errors inside their Errors
// data structure. As this function perform a breadth search, its can be a bit
// slower then the others parsers.
func (e *e121) Run() *e121 {
	for len(e.Runner) > 0 {
		switch e.Runner[0].(type) {
		case string:
			e.deleteElement()
		case []interface{}:
			e.addAll()
			e.deleteElement()
		case interface{}:
			e.processInterface()
			e.deleteElement()
		}
	}

	return e
}
