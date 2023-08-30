package gosql

type Table struct {
	name        string
	fields      []string
	constraints *Constraints
	conjunction string
	values      []any
}

func (t Table) HasField(name string) bool {
	for _, f := range t.fields {
		if f == name {
			return true
		}
	}
	return false
}

func (t *Table) Fields(fields ...string) *Table {
	for _, f := range fields {
		if !t.HasField(f) {
			t.fields = append(t.fields, f)
		}
	}
	return t
}

func (t *Table) AndWhere(field, condition string, value any) *Constraints {
	return t.constraints.And(field, condition, value)
}

func (t *Table) OrWhere(field, condition string, value any) *Constraints {
	return t.constraints.Or(field, condition, value)
}

func (t *Table) Where(field, condition string, value any) *Constraints {
	return t.AndWhere(field, condition, value)
}

func (t *Table) Values(values ...any) *Table {
	t.values = append(t.values, values...)
	return t
}

func (t *Table) Set(field string, value any) *Table {
	t.fields = append(t.fields, field)
	t.values = append(t.values, value)
	return t
}
