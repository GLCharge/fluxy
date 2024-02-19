package go_fluxy

import (
	"fmt"
	"strings"
)

type FilterOperator string
type FilterCondition string
type GroupCondition string

const (
	Eq  FilterOperator = "=="
	Neq FilterOperator = "!="
	Gt  FilterOperator = ">"
	Gte FilterOperator = ">="
	Lt  FilterOperator = "<"
	Lte FilterOperator = "<="

	And FilterCondition = "and"
	Not FilterCondition = "not"
	Or  FilterCondition = "or"

	Start GroupCondition = "start"
	End   GroupCondition = "end"
)

// Filter represents a filter with field, value, operator, and condition.
type Filter struct {
	Field     string
	Value     string
	Operator  FilterOperator
	Condition FilterCondition

	Group    GroupCondition
	Complete string
}

// FilterBuilder is a builder for constructing filters.
type FilterBuilder struct {
	filters []Filter
}

// NewFilterBuilder creates a new FilterBuilder.
func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{}
}

// AddFilter adds a filter with the specified field and value. The operator is set to equal.
// If you want to use a different operator, use the Eq(), Neq(), Gt(), Gte(), Lt(), or Lte() functions.
// For example, to add a filter with the field "foo" and value "bar" where the operator is not equal, use:
// NewFilterBuilder().AddFilter("foo", "bar").Neq()
func (fb *FilterBuilder) AddFilter(field, value string) *FilterBuilder {
	fb.filters = append(fb.filters, Filter{field, value, Eq, "", "", ""})
	return fb
}

// And adds an AND condition to the filter.
func (fb *FilterBuilder) And() *FilterBuilder {
	fb.filters = append(fb.filters, Filter{Condition: And})
	return fb
}

// Not adds an NOT condition to the filter.
func (fb *FilterBuilder) Not() *FilterBuilder {
	fb.filters = append(fb.filters, Filter{Condition: Not})
	return fb
}

// Or adds an OR condition to the filter.
func (fb *FilterBuilder) Or() *FilterBuilder {
	fb.filters = append(fb.filters, Filter{Condition: Or})
	return fb
}

// Eq sets the operator to equal.
func (fb *FilterBuilder) Eq() *FilterBuilder {
	fb.filters[len(fb.filters)-1].Operator = Eq
	return fb
}

// Neq sets the operator to not equal.
func (fb *FilterBuilder) Neq() *FilterBuilder {
	fb.filters[len(fb.filters)-1].Operator = Neq
	return fb
}

// Gt sets the operator to greater than.
func (fb *FilterBuilder) Gt() *FilterBuilder {
	fb.filters[len(fb.filters)-1].Operator = Gt
	return fb
}

// Gte sets the operator to greater than or equal.
func (fb *FilterBuilder) Gte() *FilterBuilder {
	fb.filters[len(fb.filters)-1].Operator = Gte
	return fb
}

// Lt sets the operator to less than.
func (fb *FilterBuilder) Lt() *FilterBuilder {
	fb.filters[len(fb.filters)-1].Operator = Lt
	return fb
}

// Lte sets the operator to less than or equal.
func (fb *FilterBuilder) Lte() *FilterBuilder {
	fb.filters[len(fb.filters)-1].Operator = Lte
	return fb
}

// AddCompleteFilter adds a complete filter string.
func (fb *FilterBuilder) AddCompleteFilter(filter string) *FilterBuilder {
	fb.filters = append(fb.filters, Filter{Complete: filter})
	return fb
}

// StartGroup starts a group.
func (fb *FilterBuilder) StartGroup() *FilterBuilder {
	fb.filters = append(fb.filters, Filter{Group: Start})
	return fb
}

// EndGroup ends a group.
func (fb *FilterBuilder) EndGroup() *FilterBuilder {
	fb.filters = append(fb.filters, Filter{Group: End})
	return fb
}

// Build constructs the final filter string.
func (fb *FilterBuilder) Build() string {
	var buffer strings.Builder

	for _, f := range fb.filters {
		if f.Group == Start {
			buffer.WriteString("(")
		} else if f.Group == End {
			buffer.WriteString(")")
		} else if f.Complete != "" {
			buffer.WriteString(fmt.Sprintf(`(%s)`, f.Complete))
		} else if f.Condition != "" {
			buffer.WriteString(fmt.Sprintf(" %s ", f.Condition))
		} else if f.Field != "" && f.Operator != "" {
			buffer.WriteString(fmt.Sprintf(`r.%s %s "%s"`, f.Field, f.Operator, f.Value))
		}
	}

	return buffer.String()
}
