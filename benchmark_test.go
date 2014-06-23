package jsonschema

import (
	"testing"
)

// See: https://github.com/Sembiance/cosmicrealms.com/blob/master/sandbox/benchmark-of-node-dot-js-json-validation-modules-part-2
func BenchmarkValid(b *testing.B) {
	env := RootEnv.Clone()
	schema, err := env.BuildSchema("", load_test_data("draft4/benchmark/schema4.json"))
	if err != nil {
		panic(err)
	}

	var instance interface{}
	load_test_json("draft4/benchmark/valid.json", &instance)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := schema.Validate(instance)
		if err != nil {
			b.Fatalf("error=%s", err)
		}
	}
}

func BenchmarkInvalid(b *testing.B) {
	env := RootEnv.Clone()
	schema, err := env.BuildSchema("", load_test_data("draft4/benchmark/schema4.json"))
	if err != nil {
		panic(err)
	}

	var instance interface{}
	load_test_json("draft4/benchmark/invalid.json", &instance)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := schema.Validate(instance)
		if err == nil {
			b.Fatalf("error=%s", "expected an error")
		}
	}
}
