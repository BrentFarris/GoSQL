package gosql

import (
	"strconv"
	"strings"
)

const (
	ActionSelect = "SELECT"
	ActionInsert = "INSERT"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"

	OrderDirectionAscending  = "ASC"
	OrderDirectionDescending = "DESC"
)

type Query struct {
	raw            string
	action         string
	tables         []*Table
	joins          []Join
	orderBy        string
	orderDirection string
	limit          int
	offset         int
	ignore         bool
	distinct       bool
}

func newQuery(action string) *Query {
	return &Query{
		action: action,
		tables: make([]*Table, 0),
		joins:  make([]Join, 0),
	}
}

func Select() *Query { return newQuery(ActionSelect) }
func Insert() *Query { return newQuery(ActionInsert) }
func Update() *Query { return newQuery(ActionUpdate) }
func Delete() *Query { return newQuery(ActionDelete) }
func Raw(sql string) *Query {
	q := newQuery("")
	q.raw = sql
	return q
}

func (q Query) HasConstraints() bool {
	for _, t := range q.tables {
		if len(t.constraints.constraints) > 0 {
			return true
		}
	}
	return false
}

func (q Query) IsJoin(tableName string) bool {
	for _, j := range q.joins {
		if j.rightTable == tableName {
			return true
		}
	}
	return false
}

func (q *Query) pickTable(name string) *Table {
	var t *Table = nil
	for _, t := range q.tables {
		if t.name == name {
			return t
		}
	}
	if t == nil {
		t = &Table{
			name:        name,
			fields:      make([]string, 0),
			values:      make([]any, 0),
			constraints: &Constraints{},
			conjunction: ConjunctionAnd,
		}
		q.tables = append(q.tables, t)
	}
	return t
}

func (q *Query) Table(name string) *Table { return q.pickTable(name) }
func (q *Query) From(name string) *Table  { return q.pickTable(name) }
func (q *Query) To(name string) *Table    { return q.pickTable(name) }

func (q *Query) AlsoFrom(name, conjunction string) *Table {
	t := q.Table(name)
	if t == nil {
		t = &Table{
			name:        name,
			fields:      make([]string, 0),
			constraints: &Constraints{},
			conjunction: conjunction,
		}
		q.tables = append(q.tables, t)
	}
	return t
}

func (q *Query) Join(leftTable, rightTable string) *Join {
	q.joins = append(q.joins, Join{
		leftTable:  leftTable,
		rightTable: rightTable,
	})
	return &q.joins[len(q.joins)-1]
}

func (q *Query) OrderAscending(field string) *Query {
	q.orderBy = field
	q.orderDirection = OrderDirectionAscending
	return q
}

func (q *Query) OrderDescending(field string) *Query {
	q.orderBy = field
	q.orderDirection = OrderDirectionDescending
	return q
}

func (q *Query) Ignore() *Query {
	q.ignore = true
	return q
}

func (q *Query) Distinct() *Query {
	q.distinct = true
	return q
}

func (q *Query) Limit(limit, offset int) *Query {
	q.limit = limit
	q.offset = offset
	return q
}

func (q *Query) whereString(sb *strings.Builder) []any {
	values := make([]any, 0)
	if q.HasConstraints() {
		sb.WriteString(" WHERE ")
		for i, t := range q.tables {
			if i > 0 {
				sb.WriteRune(' ')
				sb.WriteString(t.conjunction)
				sb.WriteRune(' ')
			}
			sb.WriteRune('(')
			for j, c := range t.constraints.constraints {
				if j > 0 {
					sb.WriteRune(' ')
					sb.WriteString(c.conjunction)
					sb.WriteRune(' ')
				}
				sb.WriteString(t.name)
				sb.WriteRune('.')
				sb.WriteString(c.field)
				sb.WriteString(c.condition)
				sb.WriteRune('?')
				values = append(values, c.value)
			}
			sb.WriteRune(')')
		}
	}
	return values
}

func (q *Query) selectString(sb *strings.Builder) []any {
	sb.WriteString("SELECT ")
	if q.distinct {
		sb.WriteString("DISTINCT ")
	}
	for i, t := range q.tables {
		for j, f := range t.fields {
			if i > 0 || j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(t.name)
			sb.WriteRune('.')
			sb.WriteString(f)
		}
	}
	sb.WriteString(" FROM ")
	fromCount := 0
	for _, t := range q.tables {
		if q.IsJoin(t.name) {
			continue
		} else if fromCount > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(t.name)
		fromCount++
	}
	for _, j := range q.joins {
		sb.WriteString(" LEFT JOIN ")
		sb.WriteString(j.rightTable)
		sb.WriteString(" ON ")
		sb.WriteString(j.rightTable)
		sb.WriteRune('.')
		sb.WriteString(j.rightField)
		sb.WriteRune('=')
		sb.WriteString(j.leftTable)
		sb.WriteRune('.')
		sb.WriteString(j.leftField)
	}
	values := q.whereString(sb)
	if len(q.orderBy) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(q.orderBy)
		if len(q.orderDirection) > 0 {
			sb.WriteRune(' ')
			sb.WriteString(q.orderDirection)
		}
	}
	if q.limit > 0 {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.Itoa(q.limit))
	}
	if q.offset > 0 {
		sb.WriteString(" OFFSET ")
		sb.WriteString(strconv.Itoa(q.offset))
	}
	return values
}

func (q *Query) insertString(sb *strings.Builder) []any {
	values := make([]any, 0)
	if q.ignore {
		sb.WriteString("INSERT OR IGNORE INTO ")
	} else {
		sb.WriteString("INSERT INTO ")
	}
	for i, t := range q.tables {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(t.name)
	}
	sb.WriteString(" (")
	for i, t := range q.tables {
		for j, f := range t.fields {
			if i > 0 || j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f)
		}
	}
	sb.WriteString(") VALUES (")
	for i, t := range q.tables {
		for j := range t.fields {
			if i > 0 || j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteRune('?')
		}
		values = append(values, t.values...)
	}
	sb.WriteString(")")
	return values
}

func (q *Query) updateString(sb *strings.Builder) []any {
	values := make([]any, 0)
	sb.WriteString("UPDATE ")
	for i, t := range q.tables {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(t.name)
	}
	sb.WriteString(" SET ")
	for i, t := range q.tables {
		for j, f := range t.fields {
			if i > 0 || j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(t.name)
			sb.WriteRune('.')
			sb.WriteString(f)
			sb.WriteString("=?")
			values = append(values, t.values[j])
		}
	}
	return append(values, q.whereString(sb)...)
}

func (q *Query) deleteString(sb *strings.Builder) []any {
	sb.WriteString("DELETE FROM ")
	for i, t := range q.tables {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(t.name)
	}
	return q.whereString(sb)
}

func (q *Query) Build() (string, []any) {
	if len(q.raw) > 0 {
		return q.raw, []any{}
	}
	sb := strings.Builder{}
	values := make([]any, 0)
	switch q.action {
	case ActionSelect:
		values = q.selectString(&sb)
	case ActionInsert:
		values = q.insertString(&sb)
	case ActionUpdate:
		values = q.updateString(&sb)
	case ActionDelete:
		values = q.deleteString(&sb)
	}
	return sb.String(), values
}
