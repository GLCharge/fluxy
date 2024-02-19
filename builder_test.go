package go_fluxy

import (
	"testing"
	"time"
)

func TestFluxQueryBuilder_Build(t *testing.T) {
	testCases := []struct {
		name     string
		query    *FluxQueryBuilder
		expected string
	}{
		{
			name: "basic query",
			query: NewFluxQueryBuilder().
				From("mybucket").
				Range(time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC), nil).
				Filter("r._measurement == \"mymeasurement\"").
				GroupBy("phase"),
			expected: `from(bucket: "mybucket") |> range(start: 2022-01-01T00:00:00Z) |> filter(fn: (r) => r._measurement == "mymeasurement") |> group(columns: ["phase"])`,
		},
		{
			name: "complex query",
			query: NewFluxQueryBuilder().
				From("mybucket").
				Range(time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC), nil).
				Filter("r._measurement == \"mymeasurement\"").
				Rename(map[string]string{
					"old_col": "new_col",
				}).
				Pivot("row_key", "column_key", "value_column").
				Sum(nil).
				Limit(10, nil),
			expected: `from(bucket: "mybucket") |> range(start: 2022-01-01T00:00:00Z) |> filter(fn: (r) => r._measurement == "mymeasurement") |> rename(columns: {old_col: "new_col"}) |> pivot(rowKey: ["row_key"], columnKey: ["column_key"], valueColumn: "value_column") |> sum() |> limit(n: 10)`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.query.Build()

			if actual != tc.expected {
				t.Errorf("unexpected query; expected=%s; actual=%s", tc.expected, actual)
			}
		})
	}
}
