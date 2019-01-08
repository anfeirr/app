package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsFormat(t *testing.T) {
	tests := []struct {
		scenario string
		format   string
		args     []interface{}
		expected string
		err      bool
	}{
		{
			scenario: "formatting a string",
			format:   "alert(%s)",
			args:     []interface{}{"hello"},
			expected: `alert("hello")`,
		},
		{
			scenario: "formatting an int",
			format:   "alert(%d)",
			args:     []interface{}{42},
			expected: "alert(42)",
		},
		{
			scenario: "formatting a float",
			format:   "alert(%.2f)",
			args:     []interface{}{42.21},
			expected: "alert(42.21)",
		},
		{
			scenario: "formatting a bool",
			format:   "alert(%v)",
			args:     []interface{}{false},
			expected: "alert(false)",
		},
		{
			scenario: "formatting a struct",
			format:   "alert(%v)",
			args: []interface{}{struct {
				Name string
				Age  int `json:"age"`
			}{
				Name: "Maxence",
				Age:  33,
			}},
			expected: `alert({"Name":"Maxence","age":33})`,
		},
		{
			scenario: "formatting a slice",
			format:   "alert(%v)",
			args: []interface{}{[]interface{}{
				"Maxence",
				33,
			}},
			expected: `alert(["Maxence",33])`,
		},
		{
			scenario: "formatting an array",
			format:   "alert(%v)",
			args: []interface{}{[2]interface{}{
				"Maxence",
				33,
			}},
			expected: `alert(["Maxence",33])`,
		},
		{
			scenario: "formatting a func returns an error",
			format:   "alert(%v)",
			args:     []interface{}{func() {}},
			err:      true,
		},
		{
			scenario: "formatting a non json convertible value returns an error",
			format:   "alert(%v)",
			args:     []interface{}{struct{ Fn func() }{}},
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result, err := JsFormat(test.format, test.args...)
			if test.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expected, result)
		})
	}
}
