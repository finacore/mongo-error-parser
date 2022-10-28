package meparser

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func read_error_e11000(f string) []byte {
	byteValue, _ := os.ReadFile("assets_test/e11000_test.json")

	return byteValue
}

func TestNewE11000(t *testing.T) {
	data := read_error_e11000("default")
	var root interface{}

	json.Unmarshal(data, &root)
	e11000 := newE11000(root)

	assert.NotNil(t, e11000)
	assert.Len(t, e11000.Runner, 1)
	assert.Len(t, e11000.Errors, 0)
}

func TestE11000_getWriteErrors(t *testing.T) {
	data := read_error_e11000("default")
	var root interface{}

	json.Unmarshal(data, &root)
	e11000 := newE11000(root)
	wr := e11000.getWriteErrors()

	assert.Len(t, wr, 5)
}

func TestE11000_Run(t *testing.T) {
	data := read_error_e11000("default")
	var root interface{}

	json.Unmarshal(data, &root)
	e11000 := newE11000(root)
	e11000.Run()

	assert.Len(t, e11000.Errors, 1)
}
