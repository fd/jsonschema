package jsonschema

import (
	"encoding/json"
	"reflect"
)

type Schema struct {
	Validators []Validator
	Definition map[string]interface{}
}

type Validator interface {
	Setup(x interface{}, e *Env) error
	Validate(reflect.Value, *Context)
}
