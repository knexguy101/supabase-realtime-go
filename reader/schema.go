package reader

import "reflect"

type SchemaBuilder struct {
	List map[string]reflect.Type
}

func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		List: make(map[string]reflect.Type),
	}
}

func (x *SchemaBuilder) Set(tableName string, schema interface{}) {
	x.List[tableName] = reflect.TypeOf(schema)
}
