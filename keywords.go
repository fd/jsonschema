package jsonschema

func init() {
	DefaultEnv.RegisterKeyword("type", 100, &typeValidator{})
}
