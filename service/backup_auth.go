package service

import (
	"chiller/dao"
	"chiller/http_client"

	"sync"

	
)

func BackupAuthService() {

	// get user data from auth service
	auth, err := http_client.GetAuthServiceData()
	if err != nil {
		panic(err)
	}

	// get all users in db
	bkup, err := dao.FindAllUsers()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(auth))

	// compare --> insert vs update
	for _, v := range auth {

		go func(u dao.User) {

			defer wg.Done()
			if userInBackup(u.Id, bkup) {

				dao.UpdateUser(u)
			} else {

				dao.InsertUser(u)
			}
		}(v)
	}

	wg.Wait()

}

func userInBackup(id int64, dbUsers []dao.User) bool {

	exists := false
	for _, v := range dbUsers {
		if v.Id == id {
			exists = true
		}
	}
	return exists

}

// Roles: different process because only real many-to-many
