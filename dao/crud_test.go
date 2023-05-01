package dao

import "testing"

func TestInsertRecord(t *testing.T) {

	backup := BACKUP_ALLOWANCE
	a := Allowance{Id: 2, Balance: "Encrypted balance string", UserUuid: "Encrypted User UUID"}
	err := InsertRecord(backup, "allowance", a)
	if err != nil {
		t.Log(err)
	}

}

func TestUpdateRecord(t *testing.T) {

	backup := BACKUP_ALLOWANCE
	a := Allowance{Id: 2, Balance: "Encrypted balance string", UserUuid: "Updated Encrypted User UUID"}
	m := structToMap(a)
	if len(m) != 3 {
		t.Logf("Failed to convert struct to map for: %v", a)
	}
	err := UpdateRecord(backup, "allowance", a)
	if err != nil {
		t.Log(err)
	}

}

func TestSelectRecords(t *testing.T) {

	var records []Backup
	p := " order by backup desc limit 1"
	err := SelectRecords(BACKUP_ALLOWANCE, "backup", p, &records)
	if err != nil {
		t.Log(err)
	}
	t.Logf("%v", records)

	var allowances []Allowance
	err = SelectRecords(BACKUP_ALLOWANCE, "allowance", "", &allowances)
	if err != nil {
		t.Log(err)
	}
	t.Logf("%v", allowances)

}
