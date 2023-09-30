package querybuilder

import "fmt"

func (qb *QueryBuilder) GroupBy(query string) *QueryBuilder {
	qb.groupByQuery = fmt.Sprintf("GROUP BY %s", query)
	return qb
}

func (qb *QueryBuilder) AddGroupBy(query string) *QueryBuilder {
	if !qb.HasGroupByQuery() {
		qb.GroupBy(query)
		return qb
	}

	qb.groupByQuery += fmt.Sprintf(", %s", query)
	return qb
}

func (qb *QueryBuilder) HasGroupByQuery() bool {
	return len(qb.groupByQuery) > 0
}
