package querybuilder

func (qb *QueryBuilder) Limit(n int64) *QueryBuilder {
	qb.limitQuery = replaceQueryPlaceholders("LIMIT $%d", len(qb.values)+1, 1)
	qb.values = append(qb.values, n)
	return qb
}

func (qb *QueryBuilder) Offset(n int64) *QueryBuilder {
	qb.offsetQuery = replaceQueryPlaceholders("OFFSET $%d", len(qb.values)+1, 1)
	qb.values = append(qb.values, n)
	return qb
}
