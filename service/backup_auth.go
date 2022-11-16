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

	// reconcile lookup tables: roles
	reconcileRoles(auth)

	// get all users in db
	bkUsers, err := dao.FindAllUsers()
	if err != nil {
		panic(err)
	}

	// bkRoles, err := dao.FindAllRoles()
	// if err != nil {
	// 	panic(err)
	// }

	var wg sync.WaitGroup
	wg.Add(len(auth))

	// compare --> insert vs update
	for _, v := range auth {

		go func(u dao.User) {

			defer wg.Done()
			if userInBackup(u.Id, bkUsers) {
				err = dao.UpdateUser(u)
			} else {
				err = dao.InsertUser(u)
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
// need to do update on role table first
func reconcileRoles(users []dao.User) error {

	roles := make([]dao.Role, 0)
	for _, v := range users {
		for _, ur := range v.UserRoles {
			if len(roles) == 0 || !isConsolidated(ur.Role.Id, roles) {
				roles = append(roles, ur.Role)
			}
		}
	}

	dbRoles, err := dao.FindAllRoles()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(roles) + len(dbRoles))
	for _, v := range roles {

		go func(r dao.Role) {
			defer wg.Done()

			if rolePresent(r.Id, dbRoles) && len(dbRoles) != 0 {
				err = dao.UpdateRole(r)
			} else {
				err = dao.InsertRole(r)
			}
		}(v)
	}

	// delete from backup because no longer in auth service
	for _, v := range dbRoles {

		go func(r dao.Role) {
			defer wg.Done()

			if !rolePresent(r.Id, roles) && len(roles) != 0 {
				err = dao.DeleteRole(r)
			}
		}(v)
	}

	wg.Wait()

	return err
}

func isConsolidated(id int64, rs []dao.Role) bool {
	exists := false
	for _, v := range rs {
		if v.Id == id {
			exists = true
		}
	}
	return exists
}

func rolePresent(id int64, dbRoles []dao.Role) bool {
	exists := false
	for _, v := range dbRoles {
		if v.Id == id {
			exists = true
		}
	}
	return exists
}

func updateUserRoles(rs, bkrs []dao.Role) (err error) {

	return err
}
