package document

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValueEncoder(t *testing.T) {
	tests := []struct {
		name string
		v    Value
	}{
		{"null", NewNullValue()},
		{"bool", NewBoolValue(true)},
		{"integer", NewDoubleValue(-10)},
		{"double", NewDoubleValue(-3.14)},
		{"text", NewTextValue("foo")},
		{"blob", NewBlobValue([]byte("bar"))},
		{"array", NewArrayValue(NewValueBuffer(
			NewBoolValue(true),
			NewDoubleValue(55),
			NewDoubleValue(789.58),
			NewArrayValue(NewValueBuffer(
				NewBoolValue(false),
				NewDoubleValue(100),
				NewTextValue("baz"),
			)),
			NewBlobValue([]byte("loo")),
			NewDocumentValue(
				NewFieldBuffer().
					Add("foo1", NewBoolValue(true)).
					Add("foo2", NewDoubleValue(55)).
					Add("foo3", NewArrayValue(NewValueBuffer(
						NewBoolValue(false),
						NewDoubleValue(100),
						NewTextValue("baz"),
					))),
			),
		))},
		{"document", NewDocumentValue(
			NewFieldBuffer().
				Add("foo1", NewBoolValue(true)).
				Add("foo2", NewDoubleValue(55)).
				Add("foo3", NewArrayValue(NewValueBuffer(
					NewBoolValue(false),
					NewDoubleValue(100),
					NewTextValue("baz"),
				))).
				Add("foo4", NewDocumentValue(
					NewFieldBuffer().
						Add("foo1", NewBoolValue(true)).
						Add("foo2", NewDoubleValue(55)).
						Add("foo3", NewArrayValue(NewValueBuffer(
							NewBoolValue(false),
							NewDoubleValue(100),
							NewTextValue("baz"),
						))),
				)),
		)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer

			enc := NewValueEncoder(&buf)
			err := enc.Encode(test.v)
			require.NoError(t, err)

			got, err := decodeValue(buf.Bytes())
			require.NoError(t, err)
			require.Equal(t, test.v, got)
		})
	}

	// integers are encoded as doubles to remain consistent with
	// the way unconstrained integer values are stored, including
	// integers nested inside arrays and documents.
	t.Run("integers normalize to doubles", func(t *testing.T) {
		v := NewArrayValue(NewValueBuffer(
			NewIntegerValue(1),
			NewIntegerValue(2),
			NewIntegerValue(3),
		))

		var buf bytes.Buffer
		err := NewValueEncoder(&buf).Encode(v)
		require.NoError(t, err)

		got, err := decodeValue(buf.Bytes())
		require.NoError(t, err)

		want := NewArrayValue(NewValueBuffer(
			NewDoubleValue(1),
			NewDoubleValue(2),
			NewDoubleValue(3),
		))
		require.Equal(t, want, got)
	})
}
