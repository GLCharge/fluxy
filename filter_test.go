package go_fluxy

import (
	"testing"
)

func TestFilterBuilder(t *testing.T) {
	testCases := []struct {
		name     string
		builder  func(*FilterBuilder)
		expected string
	}{
		{
			name: "single filter",
			builder: func(fb *FilterBuilder) {
				fb.AddFilter("field1", "value1")
			},
			expected: `r.field1 = "value1"`,
		},
		{
			name: "multiple filters with AND condition",
			builder: func(fb *FilterBuilder) {
				fb.AddFilter("field1", "value1").And().
					AddFilter("field2", "value2")
			},
			expected: `r.field1 = "value1" and r.field2 = "value2"`,
		},
		{
			name: "multiple filters with OR condition",
			builder: func(fb *FilterBuilder) {
				fb.AddFilter("field1", "value1").Or().
					AddFilter("field2", "value2")
			},
			expected: `r.field1 = "value1" or r.field2 = "value2"`,
		},
		{
			name: "complete filter",
			builder: func(fb *FilterBuilder) {
				fb.AddCompleteFilter(`r.field1 = "value1" and r.field2 = "value2"`)
			},
			expected: `(r.field1 = "value1" and r.field2 = "value2")`,
		},
		{
			name: "filters with mixed conditions",
			builder: func(fb *FilterBuilder) {
				fb.AddFilter("field1", "value1").And().
					AddFilter("field2", "value2").Or().
					AddFilter("field3", "value3").And().
					AddFilter("field4", "value4")
			},
			expected: `r.field1 = "value1" and r.field2 = "value2" or r.field3 = "value3" and r.field4 = "value4"`,
		},
		{
			name: "filters with different operators",
			builder: func(fb *FilterBuilder) {
				fb.AddFilter("field1", "value1").Gt().And().
					AddFilter("field2", "value2").Lt().Or().
					AddFilter("field3", "value3").Neq()
			},
			expected: `r.field1 > "value1" and r.field2 < "value2" or r.field3 != "value3"`,
		},
		{
			name: "larger filter expression",
			builder: func(fb *FilterBuilder) {
				fb.AddFilter("field1", "value1").And().
					AddFilter("field2", "value2").And().
					AddFilter("field3", "value3").Gt().Or().
					AddFilter("field4", "value4").Neq().And().
					AddFilter("field5", "value5").Lte()
			},
			expected: `r.field1 = "value1" and r.field2 = "value2" and r.field3 > "value3" or r.field4 != "value4" and r.field5 <= "value5"`,
		},
		{
			name: "Grouped filters",
			builder: func(fb *FilterBuilder) {
				fb.AddFilter("field1", "value1").Eq().
					And().
					AddFilter("field2", "value2").Gt().
					Or().
					StartGroup().
					AddFilter("field3", "value3").Lt().
					And().
					AddFilter("field4", "value4").Neq().
					EndGroup()
			},
			expected: `r.field1 = "value1" and r.field2 > "value2" or (r.field3 < "value3" and r.field4 != "value4")`,
		},
		{
			name: "Complete filters",
			builder: func(fb *FilterBuilder) {
				fb.AddCompleteFilter(`r.field1 = "value1" and r.field2 > "value2"`).
					Or().
					AddCompleteFilter(`r.field3 < "value3" and r.field4 != "value4"`)
			},
			expected: `(r.field1 = "value1" and r.field2 > "value2") or (r.field3 < "value3" and r.field4 != "value4")`,
		},
		{
			name: "Combination of all",
			builder: func(fb *FilterBuilder) {
				fb.AddFilter("field1", "value1").Eq().
					And().
					AddFilter("field2", "value2").Gt().
					Or().
					StartGroup().
					AddFilter("field3", "value3").Lt().
					And().
					AddFilter("field4", "value4").Neq().
					EndGroup().
					Or().
					AddCompleteFilter(`r.field5 = "value5" and r.field6 > "value6"`)
			},
			expected: `r.field1 = "value1" and r.field2 > "value2" or (r.field3 < "value3" and r.field4 != "value4") or (r.field5 = "value5" and r.field6 > "value6")`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fb := NewFilterBuilder()
			tc.builder(fb)
			filter := fb.Build()

			if filter != tc.expected {
				t.Errorf("Expected filter: %s, got: %s", tc.expected, filter)
			}
		})
	}
}
