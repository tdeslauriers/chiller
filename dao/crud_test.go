package dao

import "testing"

func TestInsertRecord(t *testing.T) {

	backup := BACKUP_ALLOWANCE
	cols := []string{"id, balance, user_uuid"}
	a := Allowance{Id: 1, Balance: "Encrypted balance string", User_Uuid: "Encrypted User UUID"}
	InsertRecord(backup, "allowance", cols, a.Id, a.Balance, a.User_Uuid)

}

func TestUpdateRecord(t *testing.T) {

	backup := BACKUP_ALLOWANCE
	a := Allowance{Id: 1, Balance: "Encrypted balance string", User_Uuid: "Updated Encrypted User UUID"}
	m := structToMap(a)
	if len(m) != 3 {
		t.Logf("Failed to convert struct to map for: %v", a)
	}
	err := UpdateRecord(backup, "allowance", a)
	if err != nil {
		t.Log(err)
	}

}
