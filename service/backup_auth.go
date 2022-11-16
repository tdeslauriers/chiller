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

			inBackup := false
			for _, v := range bkUsers {
				if v.Id == u.Id {
					inBackup = true
				}
			}
			if inBackup {
				err = dao.UpdateUser(u)
			} else {
				err = dao.InsertUser(u)
			}
			reconcileUserRoles(u)
		}(v)
	}

	wg.Wait()

}

// Roles: different process because only real many-to-many
// need to do update on role table first
func reconcileRoles(users []dao.User) error {

	roles := make([]dao.Role, 0)
	for _, v := range users {
		for _, ur := range v.UserRoles {

			isConsolidated := false
			for _, v := range roles {
				if v.Id == ur.Id {
					isConsolidated = true
				}
			}
			if len(roles) == 0 || !isConsolidated {
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

			exists := false
			for _, v := range dbRoles {
				if v.Id == r.Id {
					exists = true
				}
			}
			if exists && len(dbRoles) != 0 {
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

			exists := false
			for _, v := range roles {
				if v.Id == r.Id {
					exists = true
				}
			}
			if !exists && len(roles) != 0 {
				err = dao.DeleteRole(r)
			}
		}(v)
	}

	wg.Wait()

	return err
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

			exists := false // in back up data.
			for _, v := range bkur {
				if ur.Id == v.Id {
					exists = true
				}
			}
			if !exists || len(bkur) == 0 {
				err = dao.InsertUserRole(user.Id, ur)
			}
		}(v)
	}

	for _, v := range bkur {

		go func(ur dao.UrXref) {
			defer wg.Done()

			exists := false // in the auth-service json data
			for _, v := range user.UserRoles {
				if ur.Id == v.Id {
					exists = true
				}
			}

			if !exists {
				err = dao.DeleteUserRole(ur.Id)
			}
		}(v)
	}

	wg.Wait()

	return err
}
