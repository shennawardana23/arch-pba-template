package querybuilder

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	baseQuery    string
	whereQuery   string
	orderByQuery string
	limitQuery   string
	offsetQuery  string
	groupByQuery string
	tailQuery    string // for locking, returning
	values       []interface{}
}

func NewBuilder(query string, values ...interface{}) *QueryBuilder {
	if len(values) > 0 {
		query = replaceQueryPlaceholders(query, 1, len(values))
	}
	return &QueryBuilder{
		baseQuery: query,
		values:    values,
	}
}

func (qb *QueryBuilder) AppendBaseQuery(query string, values ...interface{}) *QueryBuilder {
	if len(values) > 0 {
		query = replaceQueryPlaceholders(query, len(qb.values)+1, len(values))
		qb.AppendValues(values...)
	}
	qb.baseQuery += fmt.Sprintf(" %s", query)
	return qb
}

func (qb *QueryBuilder) Build() string {
	// tailQuery should be at the end
	queries := []string{qb.baseQuery, qb.whereQuery, qb.groupByQuery, qb.orderByQuery, qb.limitQuery, qb.offsetQuery, qb.tailQuery}
	return strings.Join(queries, " ")
}

func (qb *QueryBuilder) AppendValues(values ...interface{}) {
	qb.values = append(qb.values, values...)
}

func (qb *QueryBuilder) Values() []interface{} {
	return qb.values
}

func (qb *QueryBuilder) EndsWith(query string) *QueryBuilder {
	qb.tailQuery = query
	return qb
}
