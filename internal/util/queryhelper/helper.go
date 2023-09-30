package querybuilder

import (
	"fmt"
	"reflect"
	"strings"
)

const placeholder string = "?"

// To replace query placeholders, we use '%d' in the query to represent '$1' or '$2' to make the query executable
// start represents starting count
// n represents how much it should be repeated
func replaceQueryPlaceholders(query string, start, n int) string {
	if n <= 0 {
		return query
	}

	placeholders := make([]interface{}, n)

	for i := 0; i < n; i++ {
		placeholders[i] = start + i
	}

	query = strings.ReplaceAll(query, placeholder, "$%d")
	return fmt.Sprintf(query, placeholders...)
}

// To build IN query and put placeholder inside the parentheses.
// ex: column_name IN ($1, $2, $3)
func buildInQuery(column string, start, n int) string {
	if n <= 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(column)
	sb.WriteString(" IN (")

	placeholders := make([]interface{}, n)

	for i := 0; i < n; i++ {
		sb.WriteString("$%d")
		if i < n-1 {
			sb.WriteString(", ")
		}
		placeholders[i] = start + i
	}

	sb.WriteString(")")

	return fmt.Sprintf(sb.String(), placeholders...)
}

// To build NOT IN query and put placeholder inside the parentheses.
// ex: column_name IN ($1, $2, $3)
func buildNotInQuery(column string, start, n int) string {
	if n <= 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(column)
	sb.WriteString(" NOT IN (")

	placeholders := make([]interface{}, n)

	for i := 0; i < n; i++ {
		sb.WriteString("$%d")
		if i < n-1 {
			sb.WriteString(", ")
		}
		placeholders[i] = start + i
	}

	sb.WriteString(")")

	return fmt.Sprintf(sb.String(), placeholders...)
}

func isIterable(data interface{}) bool {
	return reflect.TypeOf(data).Kind() == reflect.Slice || reflect.TypeOf(data).Kind() == reflect.Array
}

// Convert any interface to array of interface
func convertData(data interface{}) []interface{} {
	if isIterable(data) {
		values := reflect.ValueOf(data)
		result := make([]interface{}, values.Len())

		for i := 0; i < values.Len(); i++ {
			result[i] = values.Index(i).Interface()
		}

		return result
	}

	// TODO: handle Map Type

	return []interface{}{data}
}
