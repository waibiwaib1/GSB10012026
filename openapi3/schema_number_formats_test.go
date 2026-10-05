package openapi3

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefineNumberFormatCallback(t *testing.T) {
	DefineNumberFormatCallback("epsg-code", func(value float64) bool {
		return value == 4326
	})
	t.Cleanup(func() { delete(SchemaNumberFormats, "epsg-code") })

	schema := NewFloat64Schema().WithFormat("epsg-code")

	require.NoError(t, schema.VisitJSON(4326))

	err := schema.VisitJSON(1234)
	require.Error(t, err)

	var schemaErr *SchemaError
	require.ErrorAs(t, err, &schemaErr)
	require.Equal(t, "format", schemaErr.SchemaField)
}

func TestDefineNumberFormatCallbackInteger(t *testing.T) {
	DefineNumberFormatCallback("positive", func(value float64) bool {
		return value > 0
	})
	t.Cleanup(func() { delete(SchemaNumberFormats, "positive") })

	schema := NewIntegerSchema().WithFormat("positive")

	require.NoError(t, schema.VisitJSON(1, EnableFormatValidation()))
	require.Error(t, schema.VisitJSON(-1, EnableFormatValidation()))
}

func TestDefineNumberFormatCallbackDocumentValidation(t *testing.T) {
	DefineNumberFormatCallback("epsg-code", func(value float64) bool {
		return true
	})
	t.Cleanup(func() { delete(SchemaNumberFormats, "epsg-code") })

	schema := NewFloat64Schema().WithFormat("epsg-code")
	require.NoError(t, schema.Validate(context.Background(), EnableSchemaFormatValidation()))
}
