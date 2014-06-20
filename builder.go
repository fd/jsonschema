package jsonschema

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
)

type Builder interface {
	Build(pointer string, v map[string]interface{}) (*Schema, error)
}

type builder struct {
	env        *Env
	stack      []*Schema
	references map[string]*Schema
}

func newBuilder(env *Env) *builder {
	return &builder{
		env:        env,
		references: map[string]*Schema{},
	}
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
			base = b.stack[l-1].Id
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

		if base != nil {
			ref = base.ResolveReference(ref)
		}

		schema.Ref = ref
		return schema, nil
	}

	b.stack = append(b.stack, schema)
	defer func() { b.stack = b.stack[:len(b.stack)-1] }()

	validators = map[int]Validator{}
	schema.Definition = v

	for k, x := range v {
		keyword, found := b.env.keywords[k]
		if !found {
			continue
		}

		validator := reflect.New(keyword.prototype).Interface().(Validator)
		err := validator.Setup(x, b)
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
			continue
		}

		refSchema, found = b.env.schemas[ref]
		if found && refSchema != nil {
			schema.RefSchema = refSchema
			continue
		}

		return fmt.Errorf("unknown schema: %s", ref)
	}

	return nil
}
