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
	// putting here so constraints are violated if role deleted that still has xref
	reconcileRoles(auth)

	// get all users in db
	bkUsers, err := dao.FindAllUsers()
	if err != nil {
		panic(err)
	}

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
			reconcileUserRoles(u)
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

func reconcileUserRoles(user dao.User) (err error) {

	bkur, err := dao.FindUserRolesByUserId(user.Id)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(user.UserRoles) + len(bkur))
	for _, v := range user.UserRoles {

		go func(ur dao.UserRoles) {
			defer wg.Done()
			if !urInBackup(ur.Id, bkur) || len(bkur) == 0 {
				err = dao.InsertUserRole(user.Id, ur)
			}
		}(v)
	}

	for _, v := range bkur {

		go func(ur dao.UrXref) {
			defer wg.Done()
			if !urPresent(ur.Id, user.UserRoles) {
				err = dao.DeleteUserRole(ur.Id)
			}
		}(v)
	}

	wg.Wait()

	return err
}

func urInBackup(id int64, ur []dao.UrXref) bool {
	exists := false
	for _, v := range ur {
		if v.Id == id {
			exists = true
		}
	}
	return exists
}

func urPresent(id int64, ur []dao.UserRoles) bool {
	exists := false
	for _, v := range ur {
		if v.Id == id {
			exists = true
		}
	}
	return exists
}
