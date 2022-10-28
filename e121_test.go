package meparser

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var e121_error = map[string]string{
	"default":      "assets_test/e121_test.json",
	"array":        "assets_test/e121_array_test.json",
	"reason":       "assets_test/e121_reason_test.json",
	"multi_reason": "assets_test/e121_multi_reason_test.json",
	"description":  "assets_test/e121_description_test.json",
	"3_errors":     "assets_test/e121_3_test.json",
	"4_errors":     "assets_test/e121_4_test.json",
}

func read_errors_e121(f string) []byte {
	byteValue, _ := os.ReadFile(e121_error[f])

	return byteValue
}

func TestNewE121(t *testing.T) {
	data := read_errors_e121("default")
	var root interface{}

	if err := json.Unmarshal(data, &root); err != nil {
		t.Fail()
	}

	e121 := newE121(root)

	assert.NotNil(t, e121)
	assert.Len(t, e121.Runner, 1)
	assert.Len(t, e121.Errors, 0)
}

func TestE121_deleteElement(t *testing.T) {
	data := read_errors_e121("default")
	var root interface{}

	if err := json.Unmarshal(data, &root); err != nil {
		t.Fail()
	}

	e121 := newE121(root)
	e121.deleteElement()

	assert.Len(t, e121.Runner, 0)
}

func TestE121_addAll(t *testing.T) {
	data := read_errors_e121("array")
	var root interface{}

	if err := json.Unmarshal(data, &root); err != nil {
		t.Fail()
	}

	e121 := newE121(root)
	e121.addAll()

	//root + 2
	assert.Len(t, e121.Runner, 3)
}

func TestE121_addChildren(t *testing.T) {
	data := read_errors_e121("default")
	var root interface{}

	if err := json.Unmarshal(data, &root); err != nil {
		t.Fail()
	}

	e121 := newE121(root)
	e121.addChildren()

	//root + 2
	assert.Len(t, e121.Runner, 3)
}

func TestE121_addError(t *testing.T) {
	data := read_errors_e121("default")
	var root interface{}

	if err := json.Unmarshal(data, &root); err != nil {
		t.Fail()
	}

	e121 := newE121(root)
	e121.addError("test", "Hello world!")

	assert.Len(t, e121.Errors, 1)
}

func TestE121_getReason(t *testing.T) {
	tests := []struct {
		name string
		file string
		want []string
	}{
		{
			name: "one_reason",
			file: "reason",
			want: []string{"comparison failed"},
		},
		{
			name: "two_reason",
			file: "multi_reason",
			want: []string{"specified string length was not satisfied", "value was not found in enum"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := read_errors_e121(tt.file)
			var root interface{}

			if err := json.Unmarshal(data, &root); err != nil {
				t.Fail()
			}

			e121 := newE121(root)
			reason := e121.getReason()

			assert.Equal(t, reason, tt.want)
		})
	}
}

func TestE121_getErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		file string
		want []string
	}{
		{
			name: "description",
			file: "description",
			want: []string{"The value should be between 0 and 5"},
		},
		{
			name: "reason",
			file: "reason",
			want: []string{"comparison failed"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := read_errors_e121(tt.file)
			var root interface{}

			if err := json.Unmarshal(data, &root); err != nil {
				t.Fail()
			}

			e121 := newE121(root)
			err := e121.getErrorMessage(e121.Runner[0].(map[string]interface{}))

			assert.Equal(t, err, tt.want)
		})
	}
}

func TestE121_processInterface(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		hasErr bool
	}{
		{
			name:   "description",
			file:   "description",
			hasErr: true,
		},
		{
			name:   "reason",
			file:   "reason",
			hasErr: true,
		},
		{
			name:   "default",
			file:   "default",
			hasErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := read_errors_e121(tt.file)
			var root interface{}

			if err := json.Unmarshal(data, &root); err != nil {
				t.Fail()
			}

			e121 := newE121(root)
			e121.processInterface()

			if tt.hasErr {
				assert.Len(t, e121.Errors, 1)
				assert.Len(t, e121.Runner, 1)
			} else {
				assert.Len(t, e121.Errors, 0)
				assert.Len(t, e121.Runner, 3) //root + 2
			}
		})
	}
}

func TestE121_Run(t *testing.T) {
	tests := []struct {
		name string
		file string
		errs int
	}{
		{
			name: "default",
			file: "default",
			errs: 2,
		},
		{
			name: "3_errors",
			file: "3_errors",
			errs: 3,
		},
		{
			name: "4_errors",
			file: "4_errors",
			errs: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := read_errors_e121(tt.file)
			var root interface{}

			if err := json.Unmarshal(data, &root); err != nil {
				t.Fail()
			}

			e121 := newE121(root)
			e121.Run()

			assert.Len(t, e121.Runner, 0)
			assert.Len(t, e121.Errors, tt.errs)
		})
	}
}
