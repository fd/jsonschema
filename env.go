package jsonschema

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"sort"
)

type Env struct {
	// schema definition for schemas that can be interpreted by this environment
	SchemaSchema *Schema

	schemas  map[string]*Schema
	keywords map[string]keyword
}

type keyword struct {
	keyword   string
	priority  int
	prototype reflect.Type
}

func NewEnv() *Env {
	return &Env{
		schemas:  map[string]*Schema{},
		keywords: map[string]keyword{},
	}
}

func (e *Env) RegisterKeyword(key string, priority int, v Validator) {
	if _, found := e.keywords[key]; found {
		panic("keyword is already registered")
	}
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Ptr {
		panic("Validator must be a pointer")
	}
	e.keywords[key] = keyword{key, priority, rt.Elem()}
}

func (e *Env) RegisterSchema(id string, data []byte) (*Schema, error) {
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

	schema, err := e.BuildSchema(obj)
	if err != nil {
		return nil, err
	}

	if id, ok := obj["id"].(string); ok {
		schema.Id = id
	}

	if schema.Id != id {
		return nil, errors.New("schema id dit not match url")
	}

	e.schemas[id] = schema
	return schema, nil
}

func (e *Env) BuildSchema(v map[string]interface{}) (*Schema, error) {
	var (
		order      []int
		validators = map[int]Validator{}
		schema     = &Schema{}
	)

	schema.Definition = v

	for k, x := range v {
		keyword, found := e.keywords[k]
		if !found {
			continue
		}

		validator := reflect.New(keyword.prototype).Interface().(Validator)
		err := validator.Setup(x, e)
		if err != nil {
			return nil, err
		}

		order = append(order, keyword.priority)
		validators[keyword.priority] = validator
	}

	sort.Ints(order)

	for _, i := range order {
		schema.Validators = append(schema.Validators, validators[i])
	}

	return schema, nil
}
