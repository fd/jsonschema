package jsonschema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

type Env struct {
	Transport Transport

	schemas    map[string]*Schema
	validators map[string]*validator
	formats    map[string]FormatValidator
}

type Transport interface {
	Get(url string) ([]byte, error)
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

func (e *Env) Clone() *Env {

	schemas := make(map[string]*Schema, len(e.schemas))
	for k, v := range e.schemas {
		schemas[k] = v
	}

	validators := make(map[string]*validator, len(e.validators))
	for k, v := range e.validators {
		validators[k] = v
	}

	formats := make(map[string]FormatValidator, len(e.formats))
	for k, v := range e.formats {
		formats[k] = v
	}

	return &Env{
		e.Transport,
		schemas,
		validators,
		formats,
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
		obj         map[string]interface{}
		superschema string
	)

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	err := dec.Decode(&obj)
	if err != nil {
		return nil, err
	}

	if v, ok := obj["$schema"].(string); ok {
		superschema = normalizeRef(v)
	}

	if superschema == "" {
		superschema = "http://json-schema.org/draft-04/schema#"
	}

	if r, found := e.schemas[rootRef(superschema)]; found {
		if s, found := r.Subschemas[refFragment(superschema)]; found {
			err := s.Validate(obj)
			if err != nil {
				return nil, err
			}
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

func (e *Env) loadRemoteSchema(url string) (*Schema, error) {
	if e.Transport == nil {
		return nil, fmt.Errorf("remote schema loading is not enabled (missing transport)")
	}

	data, err := e.Transport.Get(url)
	if err != nil {
		return nil, err
	}

	return e.RegisterSchema("", data)
}
