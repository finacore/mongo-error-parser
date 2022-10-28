package meparser

import (
	"encoding/json"
	"testing"
)

func Benchmark_E11000_Run(b *testing.B) {

	data := read_error_e11000("default")
	var root interface{}

	if err := json.Unmarshal(data, &root); err != nil {
		b.Fail()
	}

	e11000 := newE11000(root)
	e11000.Run()
}
