package gosql

import (
	"strings"
)

const (
	ActionSelect = "SELECT"
	ActionInsert = "INSERT"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
)

type Query struct {
	raw    string
	action string
	tables []*Table
	joins  []Join
}

func newQuery(action string) Query {
	return Query{
		action: action,
		tables: make([]*Table, 0),
		joins:  make([]Join, 0),
	}
}

func Select() Query { return newQuery(ActionSelect) }
func Insert() Query { return newQuery(ActionInsert) }
func Update() Query { return newQuery(ActionUpdate) }
func Delete() Query { return newQuery(ActionDelete) }
func Raw(sql string) Query {
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
	for i, t := range q.tables {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(t.name)
	}
	for i, j := range q.joins {
		if i > 0 {
			sb.WriteRune(' ')
		}
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
	return q.whereString(sb)
}

func (q *Query) insertString(sb *strings.Builder) []any {
	values := make([]any, 0)
	sb.WriteString("INSERT INTO ")
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
