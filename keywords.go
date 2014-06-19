package jsonschema

func init() {
	DefaultEnv.RegisterKeyword("type", 100, &typeValidator{})

	// numbers
	DefaultEnv.RegisterKeyword("multipleOf", 200, &multipleOfValidator{})
	DefaultEnv.RegisterKeyword("exclusiveMaximum", 201, &exclusiveMaximumValidator{})
	DefaultEnv.RegisterKeyword("exclusiveMinimum", 202, &exclusiveMinimumValidator{})
	DefaultEnv.RegisterKeyword("maximum", 203, &maximumValidator{})
	DefaultEnv.RegisterKeyword("minimum", 204, &minimumValidator{})

	// strings
	DefaultEnv.RegisterKeyword("maxLength", 300, &maxLengthValidator{})
	DefaultEnv.RegisterKeyword("minLength", 301, &minLengthValidator{})
	DefaultEnv.RegisterKeyword("pattern", 302, &patternValidator{})

	// arrays
	DefaultEnv.RegisterKeyword("items", 400, &itemsValidator{})
	DefaultEnv.RegisterKeyword("additionalItems", 401, &additionalItemsValidator{})

	// objects
	DefaultEnv.RegisterKeyword("properties", 500, &propertiesValidator{})
}
