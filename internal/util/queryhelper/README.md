# Simplified Query Construction Proposal

## What's on this proposal?
> The main goal is to automate certain things that might reduce programmer's productivity and hopefully can remove redundant codes.

## Problem Statement
1. Codes in `repository.go` files have some line of codes that state: `if strings.Contains(query, "WHERE")` or `if strings.Contains(query, "ORDER BY")`.
2. Sometimes we have to manually insert `$1`, `$2`, or `$3` in the query, or we use `$%d` but with `len(values)`, `len(values)-1`, or `len(values)-2` as the replacer.

---

So, the solution given is creating a `Query Builder`.
Currently, the builder consists of `base_query`, `where_query`, `order_by_query`, `limit_query`, and `offset_query`.

### Scopes
#### Base Query
- Main query that appears before the main's `WHERE` clause.
```postgres
-- Example
SELECT *
FROM tables t
    JOIN that_table tt ON t.id = tt.t_id,
    (
        SELECT t1.id
        FROM another_tables t1
        WHERE t1.name = 'abc' -- this is not considered as main query
        ORDER BY t1.name DESC
        LIMIT 5
    )
```

#### Where Query
```postgres
WHERE t.id = 1
    OR CASE WHEN t.name = 'abc' THEN 'ABC' END
    AND (t.name NOT IN ('xxx', 'yyy', 'zzz')
        OR t.action IN ('aaa', 'bbb', 'ccc'))
```

#### Order By Query
```postgres
ORDER BY t.name ASC,
    CASE WHEN t.action = 'aaa' THEN 0 ELSE 1 END ASC
```

#### Limit and Offset Query
```postgres
LIMIT 10 OFFSET 20
```

### Functions
> Base Query Section
- `New(query string, values ...interface{}) *QueryBuilder`
  - Construct a base query. This function is the entrypoint.
- `AppendBaseQuery(query string, values ...interface{}) *QueryBuilder`
  - In case you need to customize the base query (like adding `join` or `subquery`) on certain conditions, this function can help you.
> Where Query Section
- `Where(query string, values ...interface{}) *QueryBuilder`
  - As the entrypoint for `WHERE` clause, you can use this if you know the where clause is there in `compile time`.
  - It will append `WHERE` clause in front.
  - ex: `Where("t.id = $%d", 1)`, result: `WHERE t.id = $1`, values: `[1]`
- `AndWhere(query string, values ...interface{}) *QueryBuilder`
  - This function is an appender for `WHERE` clause.
  - It will append ` AND` in front.
  - If the `whereQuery` is not there, it will append `WHERE` in front.
  - ex: `AndWhere("t.id = $%d", 1)`, result: ` AND WHERE t.id = $1`, values: `[1]`
- `OrWhere(query string, values ...interface{}) *QueryBuilder`
  - Same as `AndWhere` function.
  - It will use `OR` instead of `AND`.
- `In(column string, data interface{}) *QueryBuilder`
  - This function will construct `IN` query.
  - The `data` parameter will accept `interface{}` type which mean you can pass anything there.
  - Currently permitted type to be passed: `array`, `slice`, and `any singular type (int, float, bool)`
  - ex: `In("t.name", []string{"abc", "def", "ghi"})`, result: `t.name IN ($1, $2, $3)`, values: `["abc", "def", "ghi"]`
  - ex: `In("t.id", 4)`, result: `t.id IN ($1)`, values: `[4]`
- `NotIn(column string, data interface{}) *QueryBuilder`
  - Same as `In` function.
  - It will use `NOT IN` instead of `IN`.
- `WhereIn(column string, data interface{}) *QueryBuilder`
  - Same as `Where` function.
  - Since the format of `IN` query is quite complex, this function can help you.
  - The `data` parameter will accept `interface{}` type which mean you can pass anything there.
  - Currently permitted type to be passed: `array`, `slice`, and `any singular type (int, float, bool)`
  - ex: `WhereIn("t.name", []string{"abc", "def", "ghi"})`, result: `WHERE t.name IN ($1, $2, $3)`, values: `["abc", "def", "ghi"]`
  - ex: `WhereIn("t.id", 4)`, result: `WHERE t.id IN ($1)`, values: `[4]`
- `WhereNotIn(column string, data interface{}) *QueryBuilder`
  - Same as `WhereIn` function.
  - It will use `NOT IN` instead of `IN`.
- `AndWhereIn(column string, data interface{}) *QueryBuilder`
  - This function is an appender for `WHERE` clause.
  - It will construct `IN` query, and append ` AND` in front.
  - If the `WHERE` clause is not there, it will append `WHERE` in front.
  - The `data` parameter will accept `interface{}` type which mean you can pass anything there.
  - Currently permitted type to be passed: `array`, `slice`, and `any singular type (int, float, bool)`
  - ex: `AndWhereIn("t.name", []string{"abc", "def", "ghi"})`, result: ` AND t.name IN ($1, $2, $3)`, values: `["abc", "def", "ghi"]`
  - ex: `AndWhereIn("t.id", 4)`, result: ` AND t.id IN ($1)`, values: `[4]`
- `AndWhereNotIn(column string, data interface{}) *QueryBuilder`
  - Same as `AndWhereIn` function.
  - It will use `NOT IN` instead of `IN`.
- `OrWhereIn(column string, data interface{}) *QueryBuilder`
  - Same as `AndWhereIn` function.
  - It will use `OR` instead of `AND`.
- `OrWhereNotIn(column string, data interface{}) *QueryBuilder`
  - Same as `OrWhereIn` function.
  - It will use `OR` instedd of `AND`.
- `Condition(query string, values ...interface{}) *QueryBuilder`
  - It will construct condition(s) only.
  - It doesn't add anything in front.
  - ex: `Condition("(t.name = $%d OR t.id = $%d)", "abc", 1)`, result: `(t.name = $1 OR t.id = $2)`, values: `["abc", 1]`.
- `HasWhereQuery() bool`
  - Return true if length of `whereQuery` > 0.
- `And() *QueryBuilder`
  - Simply add ` AND ` in the `whereQuery`.
- `Or() *QueryBuilder`
  - Simply add ` OR ` in the `whereQuery`.
- `OpenWrap() *QueryBuilder`
  - Simply add ` (` in the `whereQuery`.
- `CloseWrap() *QueryBuilder`
  - Simply add `)` in the `whereQuery`.
> Order By Query Section
- `OrderBy(query string, t string, values ...interface{}) *QueryBuilder`
  - As the entrypoint of `ORDER BY` clause.
  - ex: `OrderBy("t.name", "asc")`, result: `ORDER BY t.name asc`
- `AddOrderBy(query string, t string, values ...interface{}) *QueryBuilder`
  - As appender for `orderBy` query.
  - If `orderBy` query is not there, it will append `ORDER BY` in front.
  - ex: `AddOrderBy("t.name", "asc")`, result: `, t.name asc`
- `HasOrderByQuery() bool`
  - Return true if length of `orderByQuery` > 0
> Pagination Query Section
- `Limit(n int64) *QueryBuilder`
- `Offset(n int64) *QueryBuilder`
> Misc Section
- `AppendValues(values ...interface{})`
  - Will append value(s) to the query builder.
  - Mostly, you need this value if you know the placeholder (`$1`, `$2`, or `$3`) of the query.
- `Build() string`
  - Returns all queries, it will join `baseQuery`, `whereQuery`, `orderByQuery`, `limitQuery`, and `offsetQuery` with single space (`" "`)
- `Values() []interface{}`
  - Returns the values that was stored in query builder.
  - It will be used for sql execution.

### How to Use
> Native Query
```postgres
SELECT t.id, t.name, t.action, tt."type", c.category
FROM tables t
    JOIN table_types tt ON t.type_id = tt.id,
    (
        SELECT tc.category
        FROM table_categories tc
        WHERE tc.name = 'xxx'
        LIMIT $1  -- 1
    ) as c
WHERE t.name IN ($2, $3, $4) -- 'abc', 'def', ghi'
    AND t.id NOT IN ($5) -- 3
    AND t."type" <> $6 -- 'xyz'
ORDER BY CASE WHEN t.action = $7 -- 'yyy'
    THEN 0 ELSE 1 ASC,
    t.name DESC
LIMIT 5 OFFSET 10
```

> With Builder
```golang
qb := querybuilder.New(`
    SELECT t.id, t.name, t.action, tt."type", c.category
    FROM tables t
        JOIN table_types tt ON t.type_id = tt.id,
        (
            SELECT tc.category
            FROM table_categories tc
            WHERE tc.name = 'xxx'
            LIMIT ?
        ) as c`, 1).
    WhereIn(`t.name`, []string{"abc", "def", "ghi"}).
    AndWhereNotIn(`t.id`, 3)
    AndWhere(`t."type" <> ?`, "xyz").
    OrderBy(`CASE WHEN t.action = ? THEN 0 ELSE 1`, "ASC", "yyy").
    AddOrderBy(`t.name`, "DESC").
    Limit(5).Offset(10)

// get the query
query := qb.Build()

// get the values
values := qb.Values()
```

### Examples
> Search Filter
```golang
// Before
if filter.Search != "" {
    splittedSearch := strings.Split(filter.Search, " ")
    for _, word := range splittedSearch {
        filterValues = append(filterValues, word)
        lenValue := len(filterValues)
        if strings.Contains(query, "WHERE") {
            query += fmt.Sprintf(` AND (j.title ILIKE '%%' || $%d || '%%' 
            OR e.name ILIKE '%%' || $%d || '%%' 
            OR ed.name ILIKE '%%' || $%d || '%%' 
            OR p.name ILIKE '%%' || $%d || '%%' 
            OR c.name ILIKE '%%' || $%d || '%%' 
            OR a.name ILIKE '%%' || $%d || '%%' OR ap.name ILIKE '%%' || $%d || '%%' 
            OR ac.name ILIKE '%%' || $%d || '%%')`, 
            lenValue, lenValue, lenValue, lenValue, 
            lenValue, lenValue, lenValue, lenValue)
        } else {
            query += fmt.Sprintf(` WHERE (j.title ILIKE '%%' || $%d || '%%' 
            OR e.name ILIKE '%%' || $%d || '%%' 
            OR ed.name ILIKE '%%' || $%d || '%%' 
            OR p.name ILIKE '%%' || $%d || '%%' 
            OR c.name ILIKE '%%' || $%d || '%%' 
            OR a.name ILIKE '%%' || $%d || '%%' 
            OR ap.name ILIKE '%%' || $%d || '%%' 
            OR ac.name ILIKE '%%' || $%d || '%%')`, 
            lenValue, lenValue, lenValue, lenValue, 
            lenValue, lenValue, lenValue, lenValue)
        }
    }
}

// After
if filter.Search != "" {
    splittedSearch := strings.Split(filter.Search, " ")
    for _, word := range splittedSearch {
        qb.AndWhere(`(j.title ILIKE '%%' || ? || '%%'`, word).
            OrWhere(`e.name ILIKE '%%' || ? || '%%'`, word).
            OrWhere(`ed.name ILIKE '%%' || ? || '%%'`, word).
            OrWhere(`p.name ILIKE '%%' || ? || '%%'`, word).
            OrWhere(`c.name ILIKE '%%' || ? || '%%'`, word).
            OrWhere(`a.name ILIKE '%%' || ? || '%%'`, word).
            OrWhere(`ap.name ILIKE '%%' || ? || '%%'`, word).
            OrWhere(`ac.name ILIKE '%%' || ? || '%%')`, word)
    }
}
```

> Sort
```golang
// Before
if len(sorts) != 0 {
    for _, o := range sorts {
        if o.Name == "distance" && !isDistanceExist {
            continue
        }
        if strings.Contains(query, "ORDER BY") {
            query += " ,"
        } else {
            query += " ORDER BY "
        }
        if o.Name == "salary_info" {
            query += fmt.Sprintf(` case when (j.id in (
                select
                    j2.id
                from
                    jobs j2
                where
                    ((j2.min_salary is not null
                        and j2.min_salary > 0)
                    or (j2.max_salary is not null
                        and j2.max_salary > 0)))) then 2 else 1 end %s`, o.Type)

        } else if o.Name == "type_of_source" {
            query += fmt.Sprintf(`j.type_of_source != %d`, state.Pod.SourceType.Direct.Id)
        } else {
            query += fmt.Sprintf(`%s %s`, o.Name, o.Type)
        }

        if o.Name == "is_verified" {
            query += ` NULLS LAST`
        }
    }
}

// After
if len(sorts) != 0 {
    for _, o := range sorts {
        if o.Name == "distance" && !hasDistance {
            continue
        }
        if o.Name == "salary_info" {
            qb.AddOrderBy(`case when (j.id in (
                select
                    j2.id
                from
                    jobs j2
                where
                    ((j2.min_salary is not null
                        and j2.min_salary > 0)
                    or (j2.max_salary is not null
                        and j2.max_salary > 0)))) then 2 else 1 end`, o.Type)
        } else if o.Name == "type_of_source" {
            qb.AddOrderBy(
                fmt.Sprintf(`j.type_of_source != %d`, state.Pod.SourceType.Direct.Id), 
                `asc`
            )
        } else {
            qb.AddOrderBy(o.Name, o.Type)
        }

        if o.Name == "is_verified" {
            qb.AddOrderBy("", `NULLS LAST`)
        }
    }
}
```

> Where Queries
```golang
if filter.FilterSalary != nil {
    if strings.Contains(query, "WHERE") {
        query += ` AND ((j.min_salary IS NOT NULL AND j.min_salary > 0) OR (j.max_salary IS NOT NULL AND j.max_salary > 0))`
    } else {
        query += ` WHERE ((j.min_salary IS NOT NULL AND j.min_salary > 0) OR (j.max_salary IS NOT NULL AND j.max_salary > 0))`
    }
}

if filter.MinSalary != 0 {
    filterValues = append(filterValues, filter.MinSalary)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND (j.min_salary >= $%d`, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE (j.min_salary >= $%d`, len(filterValues))
    }
    if filter.FilterSalary != nil {
        query += `)`
    } else {
        query += ` OR j.min_salary ISNULL)`
    }
}

if filter.MaxSalary != 0 {
    filterValues = append(filterValues, filter.MaxSalary)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND (j.max_salary <= $%d`, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE (j.max_salary <= $%d`, len(filterValues))
    }
    if filter.FilterSalary != nil {
        query += `)`
    } else {
        query += ` OR j.max_salary ISNULL)`
    }
}

if len(filter.TypeOfJob) != 0 {
    filterValues = append(filterValues, filter.TypeOfJob)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND j.type_of_employments <@ $%d`, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE j.type_of_employments <@ $%d`, len(filterValues))
    }
}

if len(filter.TypeOfShift) != 0 {
    filterValues = append(filterValues, filter.TypeOfShift)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND j.type_of_shift = ANY($%d)`, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE j.type_of_shift = ANY($%d)`, len(filterValues))
    }
}

if len(filter.TypeOfWork) != 0 {
    filterValues = append(filterValues, filter.TypeOfWork)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND j.type_of_work = ANY($%d)`, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE j.type_of_work = ANY($%d)`, len(filterValues))
    }
}

if len(filter.JobCategory) != 0 {
    filterValues = append(filterValues, filter.JobCategory)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND j.job_category_id = ANY($%d)`, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE j.job_category_id = ANY($%d)`, len(filterValues))
    }
}

if len(filter.EducationLevel) != 0 {
    filterValues = append(filterValues, filter.EducationLevel)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND ed.level = ANY($%d)`, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE ed.level = ANY($%d)`, len(filterValues))
    }
}

if filter.MinEducationLevel != 0 {
    filterValues = append(filterValues, filter.MinEducationLevel)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND ed.level <= $%d`, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE ed.level <= $%d`, len(filterValues))
    }
}

if filter.Gender != 0 {
    filterValues = append(filterValues, filter.Gender, 0)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND (j.gender = $%d OR j.gender = $%d)`, len(filterValues)-1, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE (j.gender = $%d OR j.gender = $%d)`, len(filterValues)-1, len(filterValues))
    }
}

if filter.Age != 0 {
    filterValues = append(filterValues, filter.Age, filter.Age)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND (j.min_age <= $%d OR j.min_age ISNULL) AND (j.max_age >= $%d OR j.max_age ISNULL)`, len(filterValues)-1, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE (j.min_age <= $%d OR j.min_age ISNULL) AND (j.max_age >= $%d OR j.max_age ISNULL)`, len(filterValues)-1, len(filterValues))
    }
}

if filter.Status != 0 {
    filterValues = append(filterValues, filter.Status)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND j.status = $%d`, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE j.status = $%d`, len(filterValues))
    }
}

if filter.IsDeleted != nil {
    filterValues = append(filterValues, filter.IsDeleted)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND j.is_deleted = $%d`, len(filterValues))
    } else {
        query += fmt.Sprintf(` WHERE j.is_deleted = $%d`, len(filterValues))
    }
}

// After
if filter.FilterSalary != nil {
    qb.AndWhere(`((j.min_salary IS NOT NULL AND j.min_salary > 0) 
        OR (j.max_salary IS NOT NULL AND j.max_salary > 0))`)
}

if filter.MinSalary != 0 {
    query := `(j.min_salary >= ?`

    if filter.FilterSalary != nil {
        query += ")"
    } else {
        query += `j.min_salary ISNULL)`
    }

    qb.AndWhere(query)
}

if filter.MaxSalary != 0 {
    query := `(j.max_salary <= ?`

    if filter.FilterSalary != nil {
        query += ")"
    } else {
        query += `j.max_salary ISNULL)`
    }

    qb.AndWhere(query)
}

if len(filter.TypeOfJob) != 0 {
    qb.AndWhere(`j.type_of_employments <@ ?`, filter.TypeOfJob)
}

if len(filter.TypeOfShift) != 0 {
    qb.AndWhere(`j.type_of_shift = ANY(?)`, filter.TypeOfShift)
}

if len(filter.TypeOfWork) != 0 {
    qb.AndWhere(`j.type_of_work = ANY(?)`, filter.TypeOfWork)
}

if len(filter.JobCategory) != 0 {
    qb.AndWhere(`j.job_category_id = ANY(?)`, filter.JobCategory)
}

if len(filter.EducationLevel) != 0 {
    qb.AndWhere(`ed.level = ANY(?)`, filter.EducationLevel)
}

if filter.MinEducationLevel != 0 {
    qb.AndWhere(`ed.level <= ?`, filter.MinEducationLevel)
}

if filter.Gender != 0 {
    qb.AndWhere(`(j.gender = ? OR j.gender = ?)`, filter.Gender, 0)
}

if filter.Age != 0 {
    qb.AndWhere(`(j.min_age <= ? OR j.min_age ISNULL) 
        AND (j.max_age >= ? OR j.max_age ISNULL)`, filter.Age, filter.Age)
}

if filter.Status != 0 {
    qb.AndWhere(`j.status = ?`, filter.Status)
}

if filter.IsDeleted != nil {
    qb.AndWhere(`j.is_deleted = ?`, filter.IsDeleted)
}
```

> Base Query with Condition(s)
```golang
//Before
query := `SELECT 
    j.id, 
    j.job_category_id, 
    j.title, 
    j.type_of_employments, 
    j.type_of_shift,
    j.min_salary,
    j.max_salary,
    j.status,
    j.created_at,
    j.updated_at,
    j.published_at,
    (
        select array_to_json(array_agg(row_to_json(s.*))) 
        from skills s where s.id = any(j.skills)
    ) as skills,
    (
        select array_to_json(array_agg(row_to_json(d2.*))) 
        from documents d2 where d2.id = any(j.documents)
    ) as documents,
    j.gender,
    j.min_age,
    j.max_age,
    j.expired_at,
    j.years_of_experience,
    j.type_of_work,
    j.type_of_source,
    j.is_verified,
    e.id as employer_id,
    e.name as employer_name,
    e.about as employer_about,
    e.logo_file_id as employer_logo_file_id,
    ed.id as education_id,
    ed.name as education_name,
    p.id as province_id,
    p.name as province_name,
    c.id as city_id,
    c.name as city_name,
    a.id,
    ap.id,
    ap.name,
    ac.id,
    ac.name,
    count(*) OVER() AS full_count
FROM jobs j 
JOIN employers e ON e.id = j.employer_id
LEFT JOIN provinces p ON p.id = j.province_id
LEFT JOIN cities c ON c.id = j.city_Id
LEFT JOIN districts d ON d.id = j.district_id
LEFT JOIN employer_addresses a ON a.id = j.employer_address_id
LEFT JOIN provinces ap ON ap.id = a.province_id
LEFT JOIN cities ac ON ac.id = a.city_id
JOIN educations ed ON ed.id = j.minimum_education_level`

queryWithDistance := `SELECT 
    j.id, 
    j.job_category_id, 
    j.title, 
    j.type_of_employments, 
    j.type_of_shift,
    j.min_salary,
    j.max_salary,
    j.status,
    j.created_at,
    j.updated_at,
    j.published_at,
    (select array_to_json(array_agg(row_to_json(s.*)))  from skills s where s.id = any(j.skills)) as skills,
    (select array_to_json(array_agg(row_to_json(d2.*)))  from documents d2 where d2.id = any(j.documents)) as documents,
    j.gender,
    j.min_age,
    j.max_age,
    j.expired_at,
    j.years_of_experience,
    j.type_of_work,
    j.type_of_source,
    j.is_verified,
    e.id as employer_id,
    e.name as employer_name,
    e.about as employer_about,
    e.logo_file_id as employer_logo_file_id,
    ed.id as education_id,
    ed.name as education_name,
    p.id as province_id,
    p.name as province_name,
    c.id as city_id,
    c.name as city_name,
    a.id,
    ap.id,
    ap.name,
    ac.id,
    ac.name,
    count(*) OVER() AS full_count,
    case 
        when (a.id is not null and asd.id is not null) then asd.geog <-> ST_MakePoint($1,$2)::geography
        when (a.id is not null and ad.id is not null) then ad.geog <-> ST_MakePoint($1,$2)::geography
        when (a.id is not null and ac.id is not null) then ac.geog <-> ST_MakePoint($1,$2)::geography
        when (d.id is not null) then d.geog <-> ST_MakePoint($1,$2)::geography
        else c.geog <-> ST_MakePoint($1,$2)::geography
    end as distance
FROM jobs j 
LEFT JOIN employers e ON e.id = j.employer_id
LEFT JOIN provinces p ON p.id = j.province_id
LEFT JOIN cities c ON c.id = j.city_Id
LEFT JOIN districts d ON d.id = j.district_id
LEFT JOIN employer_addresses a ON a.id = j.employer_address_id
LEFT JOIN provinces ap ON ap.id = a.province_id
LEFT JOIN cities ac ON ac.id = a.city_id
LEFT JOIN districts ad ON ad.id = a.district_id
LEFT JOIN sub_districts asd ON asd.id = a.sub_district_id
LEFT JOIN educations ed ON ed.id = j.minimum_education_level
`

var filterValues []interface{}
var result []entity.AllJob
isDistanceExist := false

// add join on applied jobs
if filter.FilterAppliedJob != 0 {
    query += ` LEFT JOIN 
    (
        select job_id, (jobseeker_id is not null) as is_applied
        from job_applicants
        where jobseeker_id = $1
    ) applied_jobs
        on applied_jobs.job_id = j.id`
    queryWithDistance += ` LEFT JOIN 
    (
        select job_id, (jobseeker_id is not null) as is_applied
        from job_applicants
        where jobseeker_id = $6
    ) applied_jobs
        on applied_jobs.job_id = j.id`
}

// filter part
if filter.Latitude != 0 && filter.Longitude != 0 {
    query = queryWithDistance
    isDistanceExist = true
    filterValues = append(filterValues, filter.Longitude, filter.Latitude, filter.Longitude, filter.Latitude, config.Get().Application.FilterJobsRadius)
    lenValue := len(filterValues)
    if strings.Contains(query, "WHERE") {
        query += fmt.Sprintf(` AND WHERE ST_DistanceSphere(
        case 
            when (a.id is not null and asd.id is not null) then asd.geog::geometry
            when (a.id is not null and ad.id is not null) then ad.geog::geometry
            when (a.id is not null and ac.id is not null) then ac.geog::geometry
            when (d.id is not null) then d.geog::geometry
            else c.geog::geometry
        end, ST_MakePoint($%d,$%d)) <= $%d`, lenValue-2, lenValue-1, lenValue)
    } else {
        query += fmt.Sprintf(` WHERE ST_DistanceSphere(
        case 
            when (a.id is not null and asd.id is not null) then asd.geog::geometry
            when (a.id is not null and ad.id is not null) then ad.geog::geometry
            when (a.id is not null and ac.id is not null) then ac.geog::geometry
            when (d.id is not null) then d.geog::geometry
            else c.geog::geometry
        end, ST_MakePoint($%d,$%d)) <= $%d`, lenValue-2, lenValue-1, lenValue)
    }
}

// After
var qb *querybuilder.QueryBuilder
hasDistance := filter.Latitude != 0 && filter.Longitude != 0

if hasDistance {
    qb = querybuilder.New(`SELECT 
            j.id, 
            j.job_category_id, 
            j.title, 
            j.type_of_employments, 
            j.type_of_shift,
            j.min_salary,
            j.max_salary,
            j.status,
            j.created_at,
            j.updated_at,
            j.published_at,
            (
                select array_to_json(array_agg(row_to_json(s.*)))
                from skills s where s.id = any(j.skills)
            ) as skills,
            (
                select array_to_json(array_agg(row_to_json(d2.*)))
                from documents d2 where d2.id = any(j.documents)
            ) as documents,
            j.gender,
            j.min_age,
            j.max_age,
            j.expired_at,
            j.years_of_experience,
            j.type_of_work,
            j.type_of_source,
            j.is_verified,
            e.id as employer_id,
            e.name as employer_name,
            e.about as employer_about,
            e.logo_file_id as employer_logo_file_id,
            ed.id as education_id,
            ed.name as education_name,
            p.id as province_id,
            p.name as province_name,
            c.id as city_id,
            c.name as city_name,
            a.id,
            ap.id,
            ap.name,
            ac.id,
            ac.name,
            count(*) OVER() AS full_count,
            case 
                when (a.id is not null and asd.id is not null) then asd.geog <-> ST_MakePoint($1,$2)::geography
                when (a.id is not null and ad.id is not null) then ad.geog <-> ST_MakePoint($1,$2)::geography
                when (a.id is not null and ac.id is not null) then ac.geog <-> ST_MakePoint($1,$2)::geography
                when (d.id is not null) then d.geog <-> ST_MakePoint($1,$2)::geography
                else c.geog <-> ST_MakePoint($1,$2)::geography
            end as distance
        FROM jobs j 
        LEFT JOIN employers e ON e.id = j.employer_id
        LEFT JOIN provinces p ON p.id = j.province_id
        LEFT JOIN cities c ON c.id = j.city_Id
        LEFT JOIN districts d ON d.id = j.district_id
        LEFT JOIN employer_addresses a ON a.id = j.employer_address_id
        LEFT JOIN provinces ap ON ap.id = a.province_id
        LEFT JOIN cities ac ON ac.id = a.city_id
        LEFT JOIN districts ad ON ad.id = a.district_id
        LEFT JOIN sub_districts asd ON asd.id = a.sub_district_id
        LEFT JOIN educations ed ON ed.id = j.minimum_education_level`)
    qb.AppendValues(filter.Latitude, filter.Longitude)

    if filter.FilterAppliedJob != 0 {
        qb.AppendBaseQuery(`LEFT JOIN 
        (
            select job_id, (jobseeker_id is not null) as is_applied
            from job_applicants
            where jobseeker_id = ?
        ) applied_jobs
            on applied_jobs.job_id = j.id`, filter.FilterAppliedJob)
    }

    qb.Where(`ST_DistanceSphere(
        case 
            when (a.id is not null and asd.id is not null) then asd.geog::geometry
            when (a.id is not null and ad.id is not null) then ad.geog::geometry
            when (a.id is not null and ac.id is not null) then ac.geog::geometry
            when (d.id is not null) then d.geog::geometry
            else c.geog::geometry
        end, ST_MakePoint(?,?)) <= ?`, filter.Longitude, filter.Latitude, config.Get().Application.FilterJobsRadius)
} else {
    qb = querybuilder.New(`SELECT 
            j.id, 
            j.job_category_id, 
            j.title, 
            j.type_of_employments, 
            j.type_of_shift,
            j.min_salary,
            j.max_salary,
            j.status,
            j.created_at,
            j.updated_at,
            j.published_at,
            (select array_to_json(array_agg(row_to_json(s.*)))  from skills s where s.id = any(j.skills)) as skills,
            (select array_to_json(array_agg(row_to_json(d2.*)))  from documents d2 where d2.id = any(j.documents)) as documents,
            j.gender,
            j.min_age,
            j.max_age,
            j.expired_at,
            j.years_of_experience,
            j.type_of_work,
            j.type_of_source,
            j.is_verified,
            e.id as employer_id,
            e.name as employer_name,
            e.about as employer_about,
            e.logo_file_id as employer_logo_file_id,
            ed.id as education_id,
            ed.name as education_name,
            p.id as province_id,
            p.name as province_name,
            c.id as city_id,
            c.name as city_name,
            a.id,
            ap.id,
            ap.name,
            ac.id,
            ac.name,
            count(*) OVER() AS full_count
        FROM jobs j 
        JOIN employers e ON e.id = j.employer_id
        LEFT JOIN provinces p ON p.id = j.province_id
        LEFT JOIN cities c ON c.id = j.city_Id
        LEFT JOIN districts d ON d.id = j.district_id
        LEFT JOIN employer_addresses a ON a.id = j.employer_address_id
        LEFT JOIN provinces ap ON ap.id = a.province_id
        LEFT JOIN cities ac ON ac.id = a.city_id
        JOIN educations ed ON ed.id = j.minimum_education_level`)

    if filter.FilterAppliedJob != 0 {
        qb.AppendBaseQuery(`LEFT JOIN 
        (
            select job_id, (jobseeker_id is not null) as is_applied
            from job_applicants
            where jobseeker_id = ?
        ) applied_jobs
            on applied_jobs.job_id = j.id`, filter.FilterAppliedJob)
    }
}
```

> Others
```golang
if filter.FilterSalary != nil {
    qb.AndWhere(`((j.min_salary IS NOT NULL AND j.min_salary > 0) 
        OR (j.max_salary IS NOT NULL AND j.max_salary > 0))`)
}

if filter.MinSalary != 0 {
    query := `(j.min_salary >= ?`

    if filter.FilterSalary != nil {
        query += ")"
    } else {
        query += `j.min_salary ISNULL)`
    }

    qb.AndWhere(query)
}

if filter.MaxSalary != 0 {
    query := `(j.max_salary <= ?`

    if filter.FilterSalary != nil {
        query += ")"
    } else {
        query += `j.max_salary ISNULL)`
    }

    qb.AndWhere(query)
}

if len(filter.TypeOfJob) != 0 {
    qb.AndWhere(`j.type_of_employments <@ ?`, filter.TypeOfJob)
}

if len(filter.TypeOfShift) != 0 {
    qb.AndWhere(`j.type_of_shift = ANY(?)`, filter.TypeOfShift)
}

if len(filter.TypeOfWork) != 0 {
    qb.AndWhere(`j.type_of_work = ANY(?)`, filter.TypeOfWork)
}

if len(filter.JobCategory) != 0 {
    qb.AndWhere(`j.job_category_id = ANY(?)`, filter.JobCategory)
}

if len(filter.EducationLevel) != 0 {
    qb.AndWhere(`ed.level = ANY(?)`, filter.EducationLevel)
}

if filter.MinEducationLevel != 0 {
    qb.AndWhere(`ed.level <= ?`, filter.MinEducationLevel)
}

if filter.Gender != 0 {
    qb.AndWhere(`(j.gender = ? OR j.gender = ?)`, filter.Gender, 0)
}

if filter.Age != 0 {
    qb.AndWhere(`(j.min_age <= ? OR j.min_age ISNULL) 
        AND (j.max_age >= ? OR j.max_age ISNULL)`, filter.Age, filter.Age)
}

if filter.Status != 0 {
    qb.AndWhere(`j.status = ?`, filter.Status)
}

if filter.IsDeleted != nil {
    qb.AndWhere(`j.is_deleted = ?`, filter.IsDeleted)
}

if filter.FilterAppliedJob != 0 {
    qb.OrderBy(`(is_applied is not null)`, `asc`)
}
```
```golang
qb = querybuilder.New(`SELECT * FROM tables t`).
    Where(`t.type ILIKE '%%' || ? || '%%'`, "NONE").
    And().
    OpenWrap().
    In(`t.id`, []int{2, 4, 5, 6, 7, 8}).
    OrWhereIn(`t.name`, "a").
    CloseWrap().
    OrderBy(`t.name`, `asc`)
```