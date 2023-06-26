package service

import (
	"chiller/dao"
	"chiller/http_client"
	"fmt"
	"log"
	"reflect"
	"sync"
)

func restoreAppTable(url, db, table string, t http_client.Bearer, records interface{}) error {

	if err := dao.SelectRecords(db, table, "", records); err != nil {
		return err
	}

	slice := reflect.ValueOf(records).Elem()
	if slice.Len() < 1 {
		log.Printf("%s table has no records which seems impossible", table)
	} else {

		var wgTable sync.WaitGroup
		wgTable.Add(slice.Len())

		for i := 0; i < slice.Len(); i++ {
			go func(index int) {
				defer wgTable.Done()

				if err := http_client.PostRecord(fmt.Sprintf("%s/%s", url, table), t, slice.Index(index).Interface()); err != nil {
					log.Fatalf("unable to restore record %v to service table: %s", slice.Index(index).Interface(), table)
				}
			}(i)
		}

		wgTable.Wait()
		log.Printf("Restore of %s table records complete", table)
	}

	return nil
}

func RestoreAuthService(t http_client.Bearer) error {

	var (
		users     []dao.User
		roles     []dao.Role
		userRoles []dao.UrXref
	)

	if err := restoreAppTable(http_client.Restore_auth_url, dao.BACKUP_AUTH, "user", t, &users); err != nil {
		return err
	}
	if err := restoreAppTable(http_client.Restore_auth_url, dao.BACKUP_AUTH, "role", t, &roles); err != nil {
		return err
	}
	if err := restoreAppTable(http_client.Restore_auth_url, dao.BACKUP_AUTH, "user_role", t, &userRoles); err != nil {
		return err
	}

	return nil
}

func RestoreAllowanceService(t http_client.Bearer) error {

	var (
		allowances         []dao.Allowance
		tasktypes          []dao.Tasktype
		tasks              []dao.Task
		tasktypeAllowances []dao.TasktypeAllowance
		taskAllowances     []dao.TaskAllowance
	)

	if err := restoreAppTable(http_client.Restore_allowance_url, dao.BACKUP_ALLOWANCE, "allowance", t, &allowances); err != nil {
		return err
	}
	if err := restoreAppTable(http_client.Restore_allowance_url, dao.BACKUP_ALLOWANCE, "tasktype", t, &tasktypes); err != nil {
		return err
	}
	if err := restoreAppTable(http_client.Restore_allowance_url, dao.BACKUP_ALLOWANCE, "task", t, &tasks); err != nil {
		return err
	}
	if err := restoreAppTable(http_client.Restore_allowance_url, dao.BACKUP_ALLOWANCE, "tasktype_allowance", t, &tasktypeAllowances); err != nil {
		return err
	}
	if err := restoreAppTable(http_client.Restore_allowance_url, dao.BACKUP_ALLOWANCE, "task_allowance", t, &taskAllowances); err != nil {
		return err
	}

	return nil
}

func RestoreGalleryService(t http_client.Bearer) error {

	var (
		albums []dao.Album
		images []dao.Image
		// albumImages []dao.AiXref
	)

	url := http_client.Restore_gallery_url
	db := dao.GALLERY_BACKUP_DB

	if err := restoreAppTable(url, db, "album", t, &albums); err != nil {
		return err
	}
	if err := restoreAppTable(url, db, "image", t, &images); err != nil {
		return err
	}
	// if err := restoreAppTable(url, db, "album_image", t, &albumImages); err != nil {
	// 	return err
	// }

	return nil
}
