package jsonschema

import (
	"reflect"
	"sort"
)

var DefaultEnv = NewEnv()

type Env struct {
	keywords map[string]keyword
}

type keyword struct {
	keyword   string
	priority  int
	prototype reflect.Type
}

func NewEnv() *Env {
	return &Env{
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
