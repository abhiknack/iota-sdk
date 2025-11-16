package generators

type Field struct {
	Name       string
	Type       string
	Validation string
}

type EntityConfig struct {
	ModuleName string
	EntityName string
	Fields     []Field
}
