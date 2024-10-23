package recordupdater

import (
	"github.com/apache/arrow/go/v17/arrow"
	"github.com/apache/arrow/go/v17/arrow/array"
	"github.com/apache/arrow/go/v17/arrow/memory"
	"github.com/cloudquery/cloudquery/plugins/transformer/json/client/schemaupdater"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

type JSONColumnsBuilder struct {
	i          int
	values     map[string][]*string
	typeSchema map[string]string
}

func NewJSONColumnsBuilder(typeSchema map[string]string, originalColumn *types.JSONArray) ColumnBuilder {
	b := &JSONColumnsBuilder{i: -1, values: make(map[string][]*string), typeSchema: typeSchema}
	for key, typ := range typeSchema {
		if typ != schemaupdater.JSONType {
			continue
		}
		b.values[key] = make([]*string, originalColumn.Len())
	}
	return b
}

func (b *JSONColumnsBuilder) AddRow(row map[string]any) {
	b.i++
	for key, typ := range b.typeSchema {
		if typ != schemaupdater.JSONType {
			continue
		}
		value, exists := row[key]
		if !exists {
			b.values[key][b.i] = nil
			continue
		}
		if v, ok := value.(string); ok {
			b.values[key][b.i] = &v
		}
	}
}

func (b *JSONColumnsBuilder) Build(key string) (arrow.Array, error) {
	if _, ok := b.values[key]; !ok {
		return nil, nil
	}
	return buildJSONColumn(b.values[key]), nil
}

func buildJSONColumn(values []*string) arrow.Array {
	bld := types.NewJSONBuilder(array.NewExtensionBuilder(memory.DefaultAllocator, types.NewJSONType()))
	for _, value := range values {
		bld.Append(value)
	}
	return bld.NewJSONArray()
}
