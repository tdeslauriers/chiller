package service

import (
	"chiller/dao"
)

func BackupRecord(uri string, table string, record interface{}) error {

	if err := dao.InsertRecord(uri, table, record); err != nil {
		if err := dao.UpdateRecord(uri, table, record); err != nil {
			return err
		}
	}
	return nil
}
