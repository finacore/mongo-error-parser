package meparser

import (
	"encoding/json"
	"testing"
)

func Benchmark_E121_Run(b *testing.B) {

	data := read_errors_e121("default")
	var root interface{}

	if err := json.Unmarshal(data, &root); err != nil {
		b.Fail()
	}

	e121 := newE121(root)
	e121.Run()
}
