package jsonschema

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
)

type Builder interface {
	Build(pointer string, v map[string]interface{}) (*Schema, error)
	GetFormatValidator(name string) FormatValidator
	GetKeyword(s string) (interface{}, bool)
}

type builder struct {
	env        *Env
	stack      []*builderStackFrame
	references map[string]*Schema
}

type builderStackFrame struct {
	schema   *Schema
	keywords map[string]bool
}

func newBuilder(env *Env) *builder {
	return &builder{
		env:        env,
		references: map[string]*Schema{},
	}
}

func (b *builder) GetFormatValidator(name string) FormatValidator {
	return b.env.formats[name]
}

func (b *builder) Build(pointer string, v map[string]interface{}) (*Schema, error) {
	var (
		order      []int
		validators map[int]Validator
		schema     = &Schema{}
		inlineId   *url.URL
		base       *url.URL
	)

	// resolve the id
	{
		var (
			id  *url.URL
			err error
		)

		if l := len(b.stack); l > 0 {
			base = b.stack[l-1].schema.Id
		}

		if x, ok := v["id"].(string); ok && x != "" {
			id, err = url.Parse(x)
			if err != nil {
				return nil, err
			}

			if base != nil {
				id = base.ResolveReference(id)
			}
		}

		{
			inlineId = &url.URL{}

			if base == nil {
				inlineId.Fragment = pointer
			} else {
				inlineId.Fragment = base.Fragment + pointer
			}

			if base != nil {
				inlineId = base.ResolveReference(inlineId)
			}
		}

		if id == nil {
			id = inlineId
		}

		schema.Id = id
		b.references[normalizeRef(schema.Id.String())] = schema
		b.references[normalizeRef(inlineId.String())] = schema
	}

	if refstr, ok := isRef(v); ok {
		ref, err := url.Parse(refstr)
		if err != nil {
			return nil, err
		}

		fragment := ref.Fragment

		if base != nil {
			ref = base.ResolveReference(ref)
		}

		ref.Fragment = fragment

		schema.Ref = ref
		return schema, nil
	}

	frame := &builderStackFrame{
		schema:   schema,
		keywords: make(map[string]bool, len(v)),
	}
	b.stack = append(b.stack, frame)
	defer func() { b.stack = b.stack[:len(b.stack)-1] }()

	validators = map[int]Validator{}
	schema.Definition = v

	var ready = map[*validator]bool{}
	for k := range v {
		validatorDef, found := b.env.validators[k]
		if !found {
			continue
		}
		if ready[validatorDef] {
			continue
		}
		ready[validatorDef] = true

		frame.keywords[k] = true

		validator := reflect.New(validatorDef.prototype).Interface().(Validator)
		err := validator.Setup(b)
		if err != nil {
			return nil, err
		}

		order = append(order, validatorDef.priority)
		validators[validatorDef.priority] = validator
	}

	for k, x := range v {
		if !frame.keywords[k] {
			if y, ok := x.(map[string]interface{}); ok && y != nil {
				_, err := b.Build("/"+escapeJSONPointer(k), y)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	sort.Ints(order)

	for _, i := range order {
		schema.Validators = append(schema.Validators, validators[i])
	}

	return schema, nil
}

func (b *builder) resolve() error {
	for _, schema := range b.references {

		if schema.Ref == nil {
			continue
		}

		if schema.RefSchema != nil {
			continue
		}

		ref := normalizeRef(schema.Ref.String())

		refSchema, found := b.references[ref]
		if found && refSchema != nil {
			schema.RefSchema = refSchema
			// fmt.Printf("%q => %q\n", normalizeRef(schema.Id.String()), normalizeRef(refSchema.Id.String()))
			continue
		}

		refSchema, found = b.env.schemas[ref]
		if found && refSchema != nil {
			schema.RefSchema = refSchema
			// fmt.Printf("%q => %q\n", normalizeRef(schema.Id.String()), normalizeRef(refSchema.Id.String()))
			continue
		}

		return fmt.Errorf("unknown schema: %s", ref)
	}

	return nil
}

func (b *builder) GetKeyword(s string) (interface{}, bool) {
	if len(b.stack) == 0 {
		return nil, false
	}
	frame := b.stack[len(b.stack)-1]
	frame.keywords[s] = true
	v, ok := frame.schema.Definition[s]
	return v, ok
}
