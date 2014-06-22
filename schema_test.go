package jsonschema

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"path"
	"testing"
)

func TestCoreIdenitySchema(t *testing.T) {
	var def map[string]interface{}
	load_test_json("core.json", &def)

	schema, err := RootEnv.BuildSchema("", load_test_data("core.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = schema.Validate(def)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
}

func TestDraft4(t *testing.T) {
	run_test_suite(t, "draft4/additionalItems.json")
	run_test_suite(t, "draft4/additionalProperties.json")
	run_test_suite(t, "draft4/allOf.json")
	run_test_suite(t, "draft4/anyOf.json")
	run_test_suite(t, "draft4/definitions.json")
	run_test_suite(t, "draft4/dependencies.json")
	run_test_suite(t, "draft4/enum.json")
	run_test_suite(t, "draft4/items.json")
	run_test_suite(t, "draft4/maxItems.json")
	run_test_suite(t, "draft4/maxLength.json")
	run_test_suite(t, "draft4/maxProperties.json")
	run_test_suite(t, "draft4/maximum.json")
	run_test_suite(t, "draft4/minItems.json")
	run_test_suite(t, "draft4/minLength.json")
	run_test_suite(t, "draft4/minProperties.json")
	run_test_suite(t, "draft4/minimum.json")
	run_test_suite(t, "draft4/multipleOf.json")
	run_test_suite(t, "draft4/not.json")
	run_test_suite(t, "draft4/oneOf.json")
	run_test_suite(t, "draft4/pattern.json")
	run_test_suite(t, "draft4/patternProperties.json")
	run_test_suite(t, "draft4/properties.json")
	run_test_suite(t, "draft4/ref.json")
	run_test_suite(t, "draft4/refRemote.json")
	run_test_suite(t, "draft4/required.json")
	run_test_suite(t, "draft4/type.json")
	run_test_suite(t, "draft4/uniqueItems.json")
}

func TestDraft4Optional(t *testing.T) {
	// run_test_suite(t, "draft4/optional/bignum.json")
	run_test_suite(t, "draft4/optional/format.json")
	run_test_suite(t, "draft4/optional/zeroTerminatedFloats.json")
}

func load_test_data(path string) []byte {
	data, err := ioutil.ReadFile("testdata/" + path)
	if err != nil {
		panic(err)
	}
	return data
}

func load_test_json(path string, v interface{}) {
	dec := json.NewDecoder(bytes.NewReader(load_test_data(path)))
	dec.UseNumber()
	err := dec.Decode(v)
	if err != nil {
		panic(err)
	}
}

func run_test_suite(t *testing.T, path string) {
	t.Logf("- %s", path)

	var suite []struct {
		Description string          `json:"description"`
		SchemaDef   json.RawMessage `json:"schema"`
		Tests       []struct {
			Description string      `json:"description"`
			Data        interface{} `json:"data"`
			Valid       bool        `json:"valid"`
		}
	}

	load_test_json(path, &suite)

	var (
		passed int
		failed int
	)

	env := RootEnv.Clone()
	env.Transport = &testTransport{}

	for _, group := range suite {
		t.Logf("  - %s:", group.Description)

		schema, err := env.BuildSchema("", group.SchemaDef)
		if err != nil {
			failed++
			t.Errorf("    error: %s", err)
			continue
		}

		for _, test := range group.Tests {
			err := schema.Validate(test.Data)
			if test.Valid && err == nil {
				passed++
				t.Logf("    \x1B[32m✓\x1B[0m %s", test.Description)
			} else if !test.Valid && err != nil {
				passed++
				t.Logf("    \x1B[32m✓\x1B[0m %s", test.Description)
			} else if test.Valid && err != nil {
				failed++
				t.Logf("    \x1B[31m✗\x1B[0m %s", test.Description)
				t.Errorf("      error: %s", err)
			} else if !test.Valid && err == nil {
				failed++
				t.Logf("    \x1B[31m✗\x1B[0m %s", test.Description)
				t.Errorf("      error: %s", "expected an error but non were generated")
			}
		}
	}

	if failed > 0 {
		t.Logf("(\x1B[32mpassed: %d\x1B[0m \x1B[31mfaild: %d\x1B[0m)", passed, failed)
	} else {
		t.Logf("(\x1B[32mpassed: %d\x1B[0m)", passed)
	}
}

type testTransport struct{}

func (t *testTransport) Get(rawurl string) ([]byte, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path.Join("testdata/draft4/remotes", u.Path))
	if err != nil {
		return nil, err
	}

	return data, nil
}
