package jsonschema

var RootEnv = NewEnv()

func init() {
	// any
	RootEnv.RegisterKeyword(&typeValidator{}, 100, "type")
	RootEnv.RegisterKeyword(&enumValidator{}, 101, "enum")
	RootEnv.RegisterKeyword(&anyOfValidator{}, 102, "anyOf")
	RootEnv.RegisterKeyword(&allOfValidator{}, 103, "allOf")
	RootEnv.RegisterKeyword(&oneOfValidator{}, 104, "oneOf")
	RootEnv.RegisterKeyword(&notValidator{}, 105, "not")
	RootEnv.RegisterKeyword(&definitionsValidator{}, 106, "definitions")
	RootEnv.RegisterKeyword(&formatValidator{}, 107, "format")

	// numbers
	RootEnv.RegisterKeyword(&multipleOfValidator{}, 200, "multipleOf")
	RootEnv.RegisterKeyword(&maximumValidator{}, 201, "maximum", "exclusiveMaximum")
	RootEnv.RegisterKeyword(&minimumValidator{}, 202, "minimum", "exclusiveMinimum")

	// strings
	RootEnv.RegisterKeyword(&maxLengthValidator{}, 300, "maxLength")
	RootEnv.RegisterKeyword(&minLengthValidator{}, 301, "minLength")
	RootEnv.RegisterKeyword(&patternValidator{}, 302, "pattern")

	// arrays
	RootEnv.RegisterKeyword(&itemsValidator{}, 400, "items", "additionalItems")
	RootEnv.RegisterKeyword(&maxItemsValidator{}, 401, "maxItems")
	RootEnv.RegisterKeyword(&minItemsValidator{}, 402, "minItems")
	RootEnv.RegisterKeyword(&uniqueItemsValidator{}, 403, "uniqueItems")

	// objects
	RootEnv.RegisterKeyword(&maxPropertiesValidator{}, 500, "maxProperties")
	RootEnv.RegisterKeyword(&minPropertiesValidator{}, 501, "minProperties")
	RootEnv.RegisterKeyword(&requiredValidator{}, 502, "required")
	RootEnv.RegisterKeyword(&propertiesValidator{}, 503, "properties", "patternProperties", "additionalProperties")
	RootEnv.RegisterKeyword(&dependenciesValidator{}, 504, "dependencies")

	RootEnv.RegisterFormat("date-time", &datetimeFormat{})
	RootEnv.RegisterFormat("email", &emailFormat{})
	RootEnv.RegisterFormat("hostname", &hostnameFormat{})
	RootEnv.RegisterFormat("ipv4", &ipv4Format{})
	RootEnv.RegisterFormat("ipv6", &ipv6Format{})
	RootEnv.RegisterFormat("regex", &regexFormat{})
	RootEnv.RegisterFormat("uri", &uriFormat{})
	RootEnv.RegisterFormat("uri-reference", &uriReferenceFormat{})

	// Set the root Schema
	schema, err := RootEnv.RegisterSchema("", draft4)
	if err != nil {
		panic(err)
	}

	err = schema.ValidateData(draft4)
	if err != nil {
		panic(err)
	}
}

var draft4 = []byte(`
	{
		"id": "http://json-schema.org/draft-04/schema#",
		"$schema": "http://json-schema.org/draft-04/schema#",
		"description": "Core schema meta-schema",
		"definitions": {
			"schemaArray": {
				"type": "array",
				"minItems": 1,
				"items": { "$ref": "#" }
			},
			"positiveInteger": {
				"type": "integer",
				"minimum": 0
			},
			"positiveIntegerDefault0": {
				"allOf": [ { "$ref": "#/definitions/positiveInteger" }, { "default": 0 } ]
			},
			"simpleTypes": {
				"enum": [ "array", "boolean", "integer", "null", "number", "object", "string" ]
			},
			"stringArray": {
				"type": "array",
				"items": { "type": "string" },
				"minItems": 1,
				"uniqueItems": true
			}
		},
		"type": "object",
		"properties": {
			"id": {
				"type": "string",
				"format": "uri-reference"
			},
			"$schema": {
				"type": "string",
				"format": "uri"
			},
			"title": {
				"type": "string"
			},
			"description": {
				"type": "string"
			},
			"default": {},
			"multipleOf": {
				"type": "number",
				"minimum": 0,
				"exclusiveMinimum": true
			},
			"maximum": {
				"type": "number"
			},
			"exclusiveMaximum": {
				"type": "boolean",
				"default": false
			},
			"minimum": {
				"type": "number"
			},
			"exclusiveMinimum": {
				"type": "boolean",
				"default": false
			},
			"maxLength": { "$ref": "#/definitions/positiveInteger" },
			"minLength": { "$ref": "#/definitions/positiveIntegerDefault0" },
			"pattern": {
				"type": "string",
				"format": "regex"
			},
			"additionalItems": {
				"anyOf": [
					{ "type": "boolean" },
					{ "$ref": "#" }
				],
				"default": {}
			},
			"items": {
				"anyOf": [
					{ "$ref": "#" },
					{ "$ref": "#/definitions/schemaArray" }
				],
				"default": {}
			},
			"maxItems": { "$ref": "#/definitions/positiveInteger" },
			"minItems": { "$ref": "#/definitions/positiveIntegerDefault0" },
			"uniqueItems": {
				"type": "boolean",
				"default": false
			},
			"maxProperties": { "$ref": "#/definitions/positiveInteger" },
			"minProperties": { "$ref": "#/definitions/positiveIntegerDefault0" },
			"required": { "$ref": "#/definitions/stringArray" },
			"additionalProperties": {
				"anyOf": [
					{ "type": "boolean" },
					{ "$ref": "#" }
				],
				"default": {}
			},
			"definitions": {
				"type": "object",
				"additionalProperties": { "$ref": "#" },
				"default": {}
			},
			"properties": {
				"type": "object",
				"additionalProperties": { "$ref": "#" },
				"default": {}
			},
			"patternProperties": {
				"type": "object",
				"additionalProperties": { "$ref": "#" },
				"default": {}
			},
			"dependencies": {
				"type": "object",
				"additionalProperties": {
					"anyOf": [
						{ "$ref": "#" },
						{ "$ref": "#/definitions/stringArray" }
					]
				}
			},
			"enum": {
				"type": "array",
				"minItems": 1,
				"uniqueItems": true
			},
			"type": {
				"anyOf": [
					{ "$ref": "#/definitions/simpleTypes" },
					{
						"type": "array",
						"items": { "$ref": "#/definitions/simpleTypes" },
						"minItems": 1,
						"uniqueItems": true
					}
				]
			},
			"allOf": { "$ref": "#/definitions/schemaArray" },
			"anyOf": { "$ref": "#/definitions/schemaArray" },
			"oneOf": { "$ref": "#/definitions/schemaArray" },
			"not": { "$ref": "#" }
		},
		"dependencies": {
			"exclusiveMaximum": [ "maximum" ],
			"exclusiveMinimum": [ "minimum" ]
		},
		"default": {}
	}
`)
