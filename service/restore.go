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
		log.Printf("%s table has not records which seems impossible", table)
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

	var users []dao.User

	if err := restoreAppTable(http_client.Restore_auth_url, dao.BACKUP_AUTH, "user", t, &users); err != nil {
		return err
	}

	return nil
}
