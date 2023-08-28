package db

import "testing"

func TestSelect(t *testing.T) {
	q := Select()
	q.From("Accounts").
		Fields("Name", "Email").
		Where("Name", ConditionEquals, "Bob").
		Or("Name", ConditionEquals, "Alice")
	q.AlsoFrom("Characters", ConjunctionAnd).
		Fields("Name", "Class").
		Where("Class", ConditionEquals, "Warrior").
		Or("Class", ConditionEquals, "Mage")
	q.Join("Accounts", "Characters").
		On("Id", "AccountId")
	str, vals := q.Build()
	if str != "SELECT Accounts.Name, Accounts.Email, Characters.Name, Characters.Class FROM Accounts, Characters LEFT JOIN Characters ON Characters.AccountId=Accounts.Id WHERE (Name=? OR Name=?) AND (Class=? OR Class=?)" {
		t.Errorf("Query string is incorrect: %s", str)
	}
	if len(vals) != 4 {
		t.Errorf("Query values is incorrect: %v", vals)
	}
	if vals[0] != "Bob" || vals[1] != "Alice" || vals[2] != "Warrior" || vals[3] != "Mage" {
		t.Errorf("Query values is incorrect: %v", vals)
	}
}

func TestInsert(t *testing.T) {
	q := Insert()
	q.To("Accounts").
		Fields("Name", "Email").
		Values("Bob", "bob@example.com")
	str, vals := q.Build()
	if str != "INSERT INTO Accounts (Name, Email) VALUES (?, ?)" {
		t.Errorf("Query string is incorrect: %s", str)
	}
	if len(vals) != 2 {
		t.Errorf("Query values is incorrect: %v", vals)
	}
	if vals[0] != "Bob" || vals[1] != "bob@example.com" {
		t.Errorf("Query values is incorrect: %v", vals)
	}
}

func TestUpdate(t *testing.T) {
	q := Update()
	q.Table("Accounts").
		Set("Name", "Bob").
		Set("Email", "bob@example.com").
		Where("Id", ConditionEquals, 1)
	str, vals := q.Build()
	if str != "UPDATE Accounts SET Name=?, Email=? WHERE (Id=?)" {
		t.Errorf("Query string is incorrect: %s", str)
	}
	if len(vals) != 3 {
		t.Errorf("Query values is incorrect: %v", vals)
	}
	if vals[0] != "Bob" || vals[1] != "bob@example.com" || vals[2] != 1 {
		t.Errorf("Query values is incorrect: %v", vals)
	}
}

func TestDelete(t *testing.T) {
	q := Delete()
	q.From("Accounts").
		Where("Id", ConditionEquals, 1)
	str, vals := q.Build()
	if str != "DELETE FROM Accounts WHERE (Id=?)" {
		t.Errorf("Query string is incorrect: %s", str)
	}
	if len(vals) != 1 {
		t.Errorf("Query values is incorrect: %v", vals)
	}
	if vals[0] != 1 {
		t.Errorf("Query values is incorrect: %v", vals)
	}
}
