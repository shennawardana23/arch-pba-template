package querybuilder

import "fmt"

func (qb *QueryBuilder) OrderBy(query string, t string, values ...interface{}) *QueryBuilder {
	if len(values) > 0 {
		query = replaceQueryPlaceholders(query, len(qb.values)+1, len(values))
		qb.AppendValues(values...)
	}
	qb.orderByQuery = fmt.Sprintf("ORDER BY %s %s", query, t)
	return qb
}

func (qb *QueryBuilder) AddOrderBy(query string, t string, values ...interface{}) *QueryBuilder {
	if !qb.HasOrderByQuery() {
		qb.OrderBy(query, t, values...)
		return qb
	}

	if len(values) > 0 {
		query = replaceQueryPlaceholders(query, len(qb.values)+1, len(values))
		qb.AppendValues(values...)
	}

	qb.orderByQuery += fmt.Sprintf(", %s %s", query, t)
	return qb
}

func (qb *QueryBuilder) HasOrderByQuery() bool {
	return len(qb.orderByQuery) > 0
}
