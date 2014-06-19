package jsonschema

import (
	"fmt"
	"reflect"
	"strings"
)

type PropertyValidationError struct {
	Property string
	Err      error
}

func (e *PropertyValidationError) Error() string {
	return fmt.Sprintf("Invalid property %q: %s", e.Property, e.Err)
}

type propertiesValidator struct {
	members map[string]*Schema
}

func (v *propertiesValidator) Setup(x interface{}, e *Env) error {
	defs, ok := x.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid 'properties' definition: %#v", x)
	}

	members := make(map[string]*Schema, len(defs))
	for k, y := range defs {
		mdef, ok := y.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid 'properties' definition: %#v", x)
		}

		schema, err := e.BuildSchema(mdef)
		if err != nil {
			return err
		}
		members[k] = schema
	}

	v.members = members
	return nil
}

func (v *propertiesValidator) Validate(x reflect.Value, ctx *Context) {
	if ctx.Type != ObjectType {
		return
	}

	for x.Kind() == reflect.Interface || x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	if x.Kind() == reflect.Map {
		for _, key := range x.MapKeys() {
			v.validate_property(key.String(), x.MapIndex(key), ctx)
		}
		return
	}

	if x.Kind() == reflect.Struct {
		t := x.Type()
		for i, l := 0, x.NumField(); i < l; i++ {
			fv := x.Field(i)
			sf := t.Field(i)

			// private
			if sf.PkgPath != "" {
				continue
			}

			tag := sf.Tag.Get("json")
			if idx := strings.IndexByte(tag, ','); idx > 0 {
				tag = tag[:idx]
			}

			// ignored
			if tag == "-" {
				continue
			}

			name := sf.Name
			if tag != "" {
				name = tag
			}

			v.validate_property(name, fv, ctx)
		}
		return
	}

	panic("unreachable")
}

func (v *propertiesValidator) validate_property(property string, x reflect.Value, ctx *Context) {
	schema, ok := v.members[property]
	if !ok {
		return
	}

	err := schema.ValidateValue(x)
	if err != nil {
		ctx.Report(&PropertyValidationError{property, err})
	}
}
