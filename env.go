package jsonschema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

type Env struct {
	// schema definition for schemas that can be interpreted by this environment
	SchemaSchema *Schema

	schemas    map[string]*Schema
	validators map[string]*validator
	formats    map[string]FormatValidator
}

type validator struct {
	keywords  []string
	priority  int
	prototype reflect.Type
}

func NewEnv() *Env {
	return &Env{
		schemas:    map[string]*Schema{},
		validators: map[string]*validator{},
		formats:    map[string]FormatValidator{},
	}
}

func (e *Env) RegisterKeyword(v Validator, priority int, key string, additionalKeys ...string) {
	keys := append(additionalKeys, key)
	sort.Strings(keys)

	for _, key := range keys {
		if _, found := e.validators[key]; found {
			panic("keyword is already registered")
		}
	}

	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Ptr {
		panic("Validator must be a pointer")
	}

	validator := &validator{keys, priority, rt.Elem()}
	for _, key := range keys {
		e.validators[key] = validator
	}
}

func (e *Env) RegisterFormat(key string, v FormatValidator) {
	if _, found := e.formats[key]; found {
		panic("format is already registered")
	}
	e.formats[key] = v
}

func (e *Env) RegisterSchema(id string, data []byte) (*Schema, error) {
	schema, err := e.BuildSchema(id, data)
	if err != nil {
		return nil, err
	}

	e.schemas[normalizeRef(schema.Id.String())] = schema
	return schema, nil
}

func (e *Env) BuildSchema(id string, data []byte) (*Schema, error) {
	var (
		obj map[string]interface{}
	)

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	err := dec.Decode(&obj)
	if err != nil {
		return nil, err
	}

	if e.SchemaSchema != nil {
		err := e.SchemaSchema.Validate(obj)
		if err != nil {
			return nil, err
		}
	}

	builder := newBuilder(e)
	schema, err := builder.Build(id, obj)
	if err != nil {
		return nil, err
	}

	err = builder.resolve()
	if err != nil {
		return nil, err
	}

	if id != "" && normalizeRef(schema.Id.String()) != id {
		return nil, fmt.Errorf("schema id dit not match url (%q != %q)", id, schema.Id)
	}

	return schema, nil
}
