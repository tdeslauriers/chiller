package service

import (
	"chiller/dao"
	"testing"
)

func TestBackupRecord(t *testing.T) {

	for i := 1; i < 100; i++ {
		a := dao.Allowance{Id: int64(i), Balance: "BALANCE", User_Uuid: "UUID"}
		err := BackupRecord(dao.BACKUP_ALLOWANCE, "allowance", a)
		if err != nil {
			t.Log(err)
		}
	}
}
