package querybuilder

import "fmt"

func (qb *QueryBuilder) Where(query string, values ...interface{}) *QueryBuilder {
	if len(values) > 0 {
		query = replaceQueryPlaceholders(query, len(qb.values)+1, len(values))
		qb.AppendValues(values...)
	}
	qb.whereQuery = fmt.Sprintf(" WHERE %s", query)
	return qb
}

func (qb *QueryBuilder) AndWhere(query string, values ...interface{}) *QueryBuilder {
	if !qb.HasWhereQuery() {
		qb.Where(query, values...)
		return qb
	}
	if len(values) > 0 {
		query = replaceQueryPlaceholders(query, len(qb.values)+1, len(values))
		qb.AppendValues(values...)
	}
	qb.whereQuery += fmt.Sprintf(" AND %s", query)
	return qb
}

func (qb *QueryBuilder) OrWhere(query string, values ...interface{}) *QueryBuilder {
	if !qb.HasWhereQuery() {
		qb.Where(query, values...)
		return qb
	}
	if len(values) > 0 {
		query = replaceQueryPlaceholders(query, len(qb.values)+1, len(values))
		qb.AppendValues(values...)
	}
	qb.whereQuery += fmt.Sprintf(" OR %s", query)
	return qb
}

func (qb *QueryBuilder) In(column string, data interface{}) *QueryBuilder {
	values := convertData(data)
	if len(values) == 0 {
		return qb
	}
	query := buildInQuery(column, len(qb.values)+1, len(values))
	qb.whereQuery += query
	qb.AppendValues(values...)
	return qb
}

func (qb *QueryBuilder) NotIn(column string, data interface{}) *QueryBuilder {
	values := convertData(data)
	if len(values) == 0 {
		return qb
	}
	query := buildNotInQuery(column, len(qb.values)+1, len(values))
	qb.whereQuery += query
	qb.AppendValues(values...)
	return qb
}

func (qb *QueryBuilder) WhereIn(column string, data interface{}) *QueryBuilder {
	values := convertData(data)
	if len(values) == 0 {
		return qb
	}
	query := buildInQuery(column, len(qb.values)+1, len(values))
	qb.whereQuery = fmt.Sprintf("WHERE %s", query)
	qb.AppendValues(values...)
	return qb
}

func (qb *QueryBuilder) WhereNotIn(column string, data interface{}) *QueryBuilder {
	values := convertData(data)
	if len(values) == 0 {
		return qb
	}
	query := buildNotInQuery(column, len(qb.values)+1, len(values))
	qb.whereQuery = fmt.Sprintf("WHERE %s", query)
	qb.AppendValues(values...)
	return qb
}

func (qb *QueryBuilder) AndWhereIn(column string, data interface{}) *QueryBuilder {
	if !qb.HasWhereQuery() {
		qb.WhereIn(column, data)
		return qb
	}
	values := convertData(data)
	if len(values) == 0 {
		return qb
	}
	query := buildInQuery(column, len(qb.values)+1, len(values))
	qb.whereQuery += fmt.Sprintf(" AND %s", query)
	qb.AppendValues(values...)
	return qb
}

func (qb *QueryBuilder) AndWhereNotIn(column string, data interface{}) *QueryBuilder {
	if !qb.HasWhereQuery() {
		qb.WhereNotIn(column, data)
		return qb
	}
	values := convertData(data)
	if len(values) == 0 {
		return qb
	}
	query := buildNotInQuery(column, len(qb.values)+1, len(values))
	qb.whereQuery += fmt.Sprintf(" AND %s", query)
	qb.AppendValues(values...)
	return qb
}

func (qb *QueryBuilder) OrWhereIn(column string, data interface{}) *QueryBuilder {
	if !qb.HasWhereQuery() {
		qb.WhereIn(column, data)
		return qb
	}
	values := convertData(data)
	if len(values) == 0 {
		return qb
	}
	query := buildInQuery(column, len(qb.values)+1, len(values))
	qb.whereQuery += fmt.Sprintf(" OR %s", query)
	qb.AppendValues(values...)
	return qb
}

func (qb *QueryBuilder) OrWhereNotIn(column string, data interface{}) *QueryBuilder {
	if !qb.HasWhereQuery() {
		qb.WhereNotIn(column, data)
		return qb
	}
	values := convertData(data)
	if len(values) == 0 {
		return qb
	}
	query := buildNotInQuery(column, len(qb.values)+1, len(values))
	qb.whereQuery += fmt.Sprintf(" OR %s", query)
	qb.AppendValues(values...)
	return qb
}

func (qb *QueryBuilder) Condition(query string, values ...interface{}) *QueryBuilder {
	if len(values) > 0 {
		query = replaceQueryPlaceholders(query, len(qb.values)+1, len(values))
		qb.AppendValues(values...)
	}
	qb.whereQuery += query
	return qb
}

func (qb *QueryBuilder) HasWhereQuery() bool {
	return len(qb.whereQuery) > 0
}

func (qb *QueryBuilder) And() *QueryBuilder {
	qb.whereQuery += " AND "
	return qb
}

func (qb *QueryBuilder) Or() *QueryBuilder {
	qb.whereQuery += " OR "
	return qb
}

func (qb *QueryBuilder) OpenWrap() *QueryBuilder {
	qb.whereQuery += " ("
	return qb
}

func (qb *QueryBuilder) CloseWrap() *QueryBuilder {
	qb.whereQuery += ")"
	return qb
}
