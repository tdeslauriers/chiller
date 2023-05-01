package service

import (
	"chiller/dao"
	"chiller/http_client"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"
)

func backupRecord(uri string, table string, record interface{}) error {

	if err := dao.InsertRecord(uri, table, record); err != nil {
		if err := dao.UpdateRecord(uri, table, record); err != nil {
			return err
		}
	}
	return nil
}

func GetLastBackup() (dao.Backup, error) {

	var last dao.Backup

	var backups []dao.Backup
	sqlParams := " ORDER BY backup DESC LIMIT 1"
	err := dao.SelectRecords(dao.BACKUP_ALLOWANCE, "backup", sqlParams, &backups)
	if err != nil {
		return last, err
	}
	last.Id = backups[0].Id
	last.Backup = backups[0].Backup

	return last, nil
}

func UpdateLastBackup() error {

	return nil
}

func backupAppTable(url, db, table string, epoch int64, t http_client.Bearer, records interface{}) error {

	endpoint := fmt.Sprintf("%s/%s/%d", url, table, epoch)
	if err := http_client.GetAppTable(endpoint, t, records); err != nil {
		return err
	}

	slice := reflect.ValueOf(records).Elem() // reflecting interface type to access its attributes/elements
	if slice.Len() < 1 {
		log.Printf("Backup is current for %s table", table)
	} else {

		var wgTable sync.WaitGroup
		wgTable.Add(slice.Len())

		for i := 0; i < slice.Len(); i++ {
			go func(index int) {
				defer wgTable.Done()

				if err := backupRecord(db, table, slice.Index(index).Interface()); err != nil {
					log.Panic("Unable to back up data record.")
				}
			}(i)

		}

		wgTable.Wait()
		log.Printf("Backup records updated for %s table", table)
	}

	return nil
}

func BackupAllowanceService(last dao.Backup, t http_client.Bearer) error {

	var (
		allowances         []dao.Allowance
		tasktypes          []dao.Tasktype
		tasktypeAllowances []dao.TasktypeAllowance
		tasks              []dao.Task
		taskAllowances     []dao.TaskAllowance
	)

	if err := backupAppTable(http_client.Backup_allowance_url, dao.BACKUP_ALLOWANCE, "allowances", last.Backup.Unix(), t, &allowances); err != nil {
		return err
	}
	if err := backupAppTable(http_client.Backup_allowance_url, dao.BACKUP_ALLOWANCE, "tasktypes", last.Backup.Unix(), t, &tasktypes); err != nil {
		return err
	}
	if err := backupAppTable(http_client.Backup_allowance_url, dao.BACKUP_ALLOWANCE, "tasktype_allowances", last.Backup.Unix(), t, &tasktypeAllowances); err != nil {
		return err
	}
	if err := backupAppTable(http_client.Backup_allowance_url, dao.BACKUP_ALLOWANCE, "tasks", last.Backup.Unix(), t, &tasks); err != nil {
		return err
	}
	if err := backupAppTable(http_client.Backup_allowance_url, dao.BACKUP_ALLOWANCE, "task_allowances", last.Backup.Unix(), t, &taskAllowances); err != nil {
		return err
	}

	// insert new 'most-recent' backup date
	b := dao.Backup{Backup: time.Now()}
	if err := dao.InsertRecord(dao.BACKUP_ALLOWANCE, "backup", b); err != nil {
		return err
	}
	log.Printf("Backup activities complete for allownace service. Most recent backup date-time: %v", b.Backup)

	return nil
}
