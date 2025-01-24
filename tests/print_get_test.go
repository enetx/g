package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestSprintfGet(t *testing.T) {
	tests := []struct {
		name     string
		template string
		args     []any
		expected string
	}{
		{
			name:     "Simple Map Access",
			template: "Value: {1.Get(key).Some}",
			args: []any{
				Map[String, String]{"key": "value"},
			},
			expected: "Value: value",
		},
		{
			name:     "Simple Map Any Access",
			template: "Value: {1.Get(key).Some}",
			args: []any{
				Map[String, any]{"key": "value"},
			},
			expected: "Value: value",
		},
		{
			name:     "Nested Map Access",
			template: "Deep Value: {1.Get(key).Some.Get(subkey).Some}",
			args: []any{
				Map[string, Map[String, string]]{
					"key": {"subkey": "deepvalue"},
				},
			},
			expected: "Deep Value: deepvalue",
		},
		{
			name:     "Map with Float Keys",
			template: "Float Key: {1.Get(3_14).Some}",
			args: []any{
				Map[Float, string]{3.14: "pi"},
			},
			expected: "Float Key: pi",
		},
		{
			name:     "Slice Index Access",
			template: "Index 1: {1.Get(1)}",
			args: []any{
				Slice[string]{"first", "second", "third"},
			},
			expected: "Index 1: second",
		},
		{
			name:     "Nested Slice Access",
			template: "Nested Index: {1.Get(1).Get(0)}",
			args: []any{
				Slice[Slice[Int]]{{100, 200}, {300, 400}},
			},
			expected: "Nested Index: 300",
		},
		{
			name:     "Map with Int Keys",
			template: "Int Key: {1.Get(42).Some}",
			args: []any{
				Map[int, string]{42: "intvalue"},
			},
			expected: "Int Key: intvalue",
		},
		{
			name:     "Boolean Key Map",
			template: "Bool Key: {1.Get(true).Some}",
			args: []any{
				Map[bool, string]{true: "boolvalue"},
			},
			expected: "Bool Key: boolvalue",
		},
		{
			name:     "Full Complexity",
			template: "Access: {1.Get(map).Some.Get(slice).Some.Get(1).Some.Get(0).Get(field).Some}",
			args: []any{
				Map[String, Map[string, Map[String, Slice[Map[string, string]]]]]{
					"map": {
						"slice": {
							"1": Slice[Map[string, string]]{
								{"field": "subfieldvalue"},
							},
						},
					},
				},
			},
			expected: "Access: subfieldvalue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sprintf(tt.template, tt.args...)
			if result != String(tt.expected) {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestSprintfGetNamed(t *testing.T) {
	tests := []struct {
		name     string
		template string
		named    Named
		expected string
	}{
		{
			name:     "Simple Map Access",
			template: "Value: {map.Get(key).Some}",
			named:    Named{"map": Map[String, String]{"key": "value"}},
			expected: "Value: value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sprintf(tt.template, tt.named)
			if result != String(tt.expected) {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}
