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
			template: "Value: {1.$get(key)}",
			args: []any{
				map[String]String{"key": "value"},
			},
			expected: "Value: value",
		},
		{
			name:     "Simple Map Any Access",
			template: "Value: {1.$get(key)}",
			args: []any{
				map[String]any{"key": "value"},
			},
			expected: "Value: value",
		},
		{
			name:     "Nested Map Access",
			template: "Deep Value: {1.$get(key.subkey)}",
			args: []any{
				map[string]map[String]string{
					"key": {"subkey": "deepvalue"},
				},
			},
			expected: "Deep Value: deepvalue",
		},
		{
			name:     "Map with Float Keys",
			template: "Float Key: {1.$get(3_14)}",
			args: []any{
				map[Float]string{3.14: "pi"},
			},
			expected: "Float Key: pi",
		},
		{
			name:     "Slice Index Access",
			template: "Index 1: {1.$get(1)}",
			args: []any{
				[]string{"first", "second", "third"},
			},
			expected: "Index 1: second",
		},
		{
			name:     "Nested Slice Access",
			template: "Nested Index: {1.$get(1.0)}",
			args: []any{
				[][]Int{{100, 200}, {300, 400}},
			},
			expected: "Nested Index: 300",
		},
		{
			name:     "Struct Field Access",
			template: "Struct Field: {1.$get(Field)}",
			args: []any{
				struct {
					Field string
				}{Field: "fieldvalue"},
			},
			expected: "Struct Field: fieldvalue",
		},
		{
			name:     "Complex Struct Field Access",
			template: "Complex Struct: {1.$get(SubStruct.InnerField)}",
			args: []any{
				struct {
					SubStruct struct {
						InnerField String
					}
				}{SubStruct: struct {
					InnerField String
				}{InnerField: "inner"}},
			},
			expected: "Complex Struct: inner",
		},
		{
			name:     "Map with Int Keys",
			template: "Int Key: {1.$get(42)}",
			args: []any{
				map[int]string{42: "intvalue"},
			},
			expected: "Int Key: intvalue",
		},
		{
			name:     "Boolean Key Map",
			template: "Bool Key: {1.$get(true)}",
			args: []any{
				map[bool]string{true: "boolvalue"},
			},
			expected: "Bool Key: boolvalue",
		},
		{
			name:     "Full Complexity",
			template: "Access: {1.$get(map.slice.1.0.field)}",
			args: []any{
				map[String]map[string]map[String][]map[string]string{
					"map": {
						"slice": {
							"1": []map[string]string{
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
			template: "Value: {map.$get(key)}",
			named:    Named{"map": map[String]String{"key": "value"}},
			expected: "Value: value",
		},
		{
			name:     "Full Complexity",
			template: "Access: {complex.$get(map.slice.1.0.field)} {struct.$get(SubStruct.InnerField)}",
			named: Named{
				"complex": map[String]map[string]map[String][]map[string]string{
					"map": {
						"slice": {
							"1": []map[string]string{
								{"field": "subfieldvalue"},
							},
						},
					},
				},
				"struct": struct {
					SubStruct struct {
						InnerField String
					}
				}{SubStruct: struct {
					InnerField String
				}{InnerField: "inner"}},
			},
			expected: "Access: subfieldvalue inner",
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
