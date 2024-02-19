package go_fluxy

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

type FluxQueryBuilder struct {
	queryParts []string
}

// NewFluxQueryBuilder creates a new FluxQueryBuilder
func NewFluxQueryBuilder() *FluxQueryBuilder {
	return &FluxQueryBuilder{}
}

func (qb *FluxQueryBuilder) From(bucket string) *FluxQueryBuilder {
	qb.queryParts = append(qb.queryParts, fmt.Sprintf(`from(bucket: "%s")`, bucket))
	return qb
}

func (qb *FluxQueryBuilder) Range(start time.Time, stop *time.Time) *FluxQueryBuilder {
	rangeQuery := []string{fmt.Sprintf("start: %s", start.Format(time.RFC3339))}

	if stop != nil {
		stopTime := *stop
		rangeQuery = append(rangeQuery, fmt.Sprintf("stop: %s", stopTime.Format(time.RFC3339)))
	}

	qb.queryParts = append(qb.queryParts, fmt.Sprintf("range(%s)", strings.Join(rangeQuery, ", ")))
	return qb
}

func (qb *FluxQueryBuilder) Filter(fn string) *FluxQueryBuilder {
	qb.queryParts = append(qb.queryParts, fmt.Sprintf("filter(fn: (r) => %s)", fn))
	return qb
}

func (qb *FluxQueryBuilder) GroupBy(columns ...string) *FluxQueryBuilder {
	columns = lo.Map(columns, func(item string, idx int) string {
		return fmt.Sprintf(`"%s"`, item)
	})
	qb.queryParts = append(qb.queryParts, fmt.Sprintf(`group(columns: [%s])`, strings.Join(columns, ", ")))
	return qb
}

func (qb *FluxQueryBuilder) Map(fn string) *FluxQueryBuilder {
	qb.queryParts = append(qb.queryParts, fmt.Sprintf("map(fn: (r) => ({%s}))", fn))
	return qb
}

func (qb *FluxQueryBuilder) TruncateTimeColumn(seconds int) *FluxQueryBuilder {
	qb.queryParts = append(qb.queryParts, fmt.Sprintf(`truncateTimeColumn(unit: %ds)`, seconds))
	return qb
}

func (qb *FluxQueryBuilder) Select(columns ...string) *FluxQueryBuilder {
	columns = lo.Map(columns, func(item string, idx int) string {
		return fmt.Sprintf(`"%s"`, item)
	})
	qb.queryParts = append(qb.queryParts, fmt.Sprintf("select(columns: [%s], )", strings.Join(columns, ", ")))
	return qb
}

func (qb *FluxQueryBuilder) Keep(columns ...string) *FluxQueryBuilder {
	columns = lo.Map(columns, func(item string, idx int) string {
		return fmt.Sprintf(`"%s"`, item)
	})
	qb.queryParts = append(qb.queryParts, fmt.Sprintf("keep(columns: [%s])", strings.Join(columns, ", ")))
	return qb
}

func (qb *FluxQueryBuilder) Drop(columns ...string) *FluxQueryBuilder {
	columns = lo.Map(columns, func(item string, idx int) string {
		return fmt.Sprintf(`"%s"`, item)
	})
	qb.queryParts = append(qb.queryParts, fmt.Sprintf(`drop(columns: [%s])`, strings.Join(columns, ", ")))
	return qb
}

func (qb *FluxQueryBuilder) Rename(mapping map[string]string) *FluxQueryBuilder {
	var renameArgs []string
	for k, v := range mapping {
		renameArgs = append(renameArgs, fmt.Sprintf(`%s: "%s"`, k, v))
	}
	qb.queryParts = append(qb.queryParts, fmt.Sprintf("rename(columns: {%s})", strings.Join(renameArgs, ", ")))
	return qb
}

func (qb *FluxQueryBuilder) Sort(descending bool, columns ...string) *FluxQueryBuilder {
	columns = lo.Map(columns, func(item string, idx int) string {
		return fmt.Sprintf(`"%s"`, item)
	})
	qb.queryParts = append(qb.queryParts, fmt.Sprintf(`sort(columns: [%s], desc: %s)`, strings.Join(columns, ", "), strconv.FormatBool(descending)))
	return qb
}

func (qb *FluxQueryBuilder) Pivot(rowKey, columnKey string, valueColumn string) *FluxQueryBuilder {
	qb.queryParts = append(qb.queryParts, fmt.Sprintf(`pivot(rowKey: ["%s"], columnKey: ["%s"], valueColumn: "%s")`, rowKey, columnKey, valueColumn))
	return qb
}

func (qb *FluxQueryBuilder) Sum(column *string) *FluxQueryBuilder {
	sumQuery := []string{}
	if column != nil {
		sumQuery = append(sumQuery, fmt.Sprintf(`column: "%s"`, *column))
	}
	qb.queryParts = append(qb.queryParts, fmt.Sprintf("sum(%s)", strings.Join(sumQuery, ", ")))
	return qb
}

func (qb *FluxQueryBuilder) Mean() *FluxQueryBuilder {
	qb.queryParts = append(qb.queryParts, "mean()")
	return qb
}

func (qb *FluxQueryBuilder) Median() *FluxQueryBuilder {
	qb.queryParts = append(qb.queryParts, "median()")
	return qb
}

func (qb *FluxQueryBuilder) Min() *FluxQueryBuilder {
	qb.queryParts = append(qb.queryParts, "min()")
	return qb
}

func (qb *FluxQueryBuilder) Max() *FluxQueryBuilder {
	qb.queryParts = append(qb.queryParts, "max()")
	return qb
}

// Limit limits the number of records returned
func (qb *FluxQueryBuilder) Limit(n int, offset *int) *FluxQueryBuilder {
	limitQuery := []string{fmt.Sprintf("n: %d", n)}
	if offset != nil {
		limitQuery = append(limitQuery, fmt.Sprintf("offset: %d", *offset))
	}

	qb.queryParts = append(qb.queryParts, fmt.Sprintf("limit(%s)", strings.Join(limitQuery, ", ")))
	return qb
}

func (qb *FluxQueryBuilder) Elapsed(columnName *string, timeColumn *string, unit *string) *FluxQueryBuilder {
	elapsedQuery := []string{}
	if columnName != nil {
		elapsedQuery = append(elapsedQuery, fmt.Sprintf(`columnName: "%s"`, *columnName))
	}
	if timeColumn != nil {
		elapsedQuery = append(elapsedQuery, fmt.Sprintf(`timeColumn: "%s"`, *timeColumn))
	}
	if unit != nil {
		elapsedQuery = append(elapsedQuery, fmt.Sprintf(`unit: %s`, *unit))
	}
	qb.queryParts = append(qb.queryParts, fmt.Sprintf("elapsed(%s)", strings.Join(elapsedQuery, ", ")))
	return qb
}

// Build returns the built flux query
func (qb *FluxQueryBuilder) Build() string {
	return strings.Join(qb.queryParts, " |> ")
}
