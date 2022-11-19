package service

import (
	"chiller/dao"
	"chiller/http_client"
	"log"

	"sync"
)

func BackupAuthService() {

	// get user data from users service
	users, err := http_client.GetAuthServiceData()
	if err != nil {
		panic(err)
	}
	urs, rs := reconstructRoleTables(users)
	uas, as := reconstructAddressTables(users)
	ups, ps := reconstructPhoneTables(users)

	// BACKUP RECORD INSERTIONS/UPDATES
	var wgTables sync.WaitGroup
	wgTables.Add(len(users) + len(rs) + len(as) + len(ps))

	// insert or update user backup records
	for _, v := range users {
		go func(u dao.User) {
			defer wgTables.Done()

			if err := dao.InsertUser(u); err != nil {
				if err := dao.UpdateUser(u); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}

	// insert or update role backup records
	for _, v := range rs {
		go func(r dao.Role) {
			defer wgTables.Done()

			if err := dao.InsertRole(r); err != nil {
				if err := dao.UpdateRole(r); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}

	// insert or update address backup records
	for _, v := range as {
		go func(a dao.Address) {
			defer wgTables.Done()

			if err := dao.InsertAddress(a); err != nil {
				if err := dao.UpdateAddress(a); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}

	// insert or update phone backup records
	for _, v := range ps {
		go func(p dao.Phone) {
			defer wgTables.Done()

			if err := dao.InsertPhone(p); err != nil {
				if err := dao.UpdatePhone(p); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}

	wgTables.Wait()

	var wgXref sync.WaitGroup
	wgXref.Add(len(urs) + len(uas) + len(ups))

	// insert user_role backup(xref) records
	for _, v := range urs {
		go func(ur dao.UrXref) {
			defer wgXref.Done()

			r := dao.XrefRecord[dao.UrXref]{Id: ur.Id, Fk_1: ur.User_id, Fk_2: ur.Role_id}
			dao.InsertXrefRecord(r, dao.INSERT_UR) // dumping error

		}(v)
	}

	// insert user_address(xref) backup records
	for _, v := range uas {
		go func(ua dao.UaXref) {
			defer wgXref.Done()

			r := dao.XrefRecord[dao.UaXref]{Id: ua.Id, Fk_1: ua.User_id, Fk_2: ua.Address_id}
			dao.InsertXrefRecord(r, dao.INSERT_UA)
		}(v)
	}

	// insert user_phone(xref) backup records
	for _, v := range ups {
		go func(up dao.UpXref) {
			defer wgXref.Done()

			r := dao.XrefRecord[dao.UpXref]{Id: up.Id, Fk_1: up.User_id, Fk_2: up.Phone_id}
			dao.InsertXrefRecord(r, dao.INSERT_UP)

		}(v)
	}

	wgXref.Wait()

	// DELETIONS
	// delete backup records no longer present in auth-service data
	// must do deletions in xrefs first to prevent constrain violations
	var wgDelXref sync.WaitGroup
	wgDelXref.Add(3)

	go func(userRoles []dao.UrXref) {
		defer wgDelXref.Done()
		deleteUserRolesFromBackup(userRoles)
	}(urs)

	go func(userAddresses []dao.UaXref) {
		defer wgDelXref.Done()
		deleteUserAddressesFromBackup(userAddresses)
	}(uas)

	go func(userPhones []dao.UpXref) {
		defer wgDelXref.Done()
		deleteUserPhonesFromBackup(userPhones)
	}(ups)

	wgDelXref.Wait()

	// primary table deletions
	var wgDelTables sync.WaitGroup
	wgDelTables.Add(4)

	go func(us []dao.User) {
		defer wgDelTables.Done()
		deleteUsersFromBackup(us)
	}(users)

	go func(roles []dao.Role) {
		defer wgDelTables.Done()
		deleteRolesFromBackup(roles)
	}(rs)

	go func(addresses []dao.Address) {
		defer wgDelTables.Done()
		deleteAddressesFromBackup(addresses)
	}(as)

	go func(phones []dao.Phone) {
		defer wgDelTables.Done()
		deletePhonesFromBackUp(phones)
	}(ps)

	wgDelTables.Wait()
	log.Print("Completed backup activites of auth-service.")
}

func deletePhonesFromBackUp(phones []dao.Phone) {
	bkPhones, err := dao.FindAllPhones()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(bkPhones))
	for _, v := range bkPhones {
		go func(p dao.Phone) {
			defer wg.Done()

			exists := false
			for _, phone := range phones {
				if p.Id == phone.Id {
					exists = true
				}
			}
			if !exists {
				record := dao.Record[dao.Phone]{Id: p.Id}
				if err := dao.DeleteRecord(record, dao.DELETE_PHONE); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}
	wg.Wait()
}

func deleteAddressesFromBackup(addresses []dao.Address) {
	bkAddresses, err := dao.FindAllAddresses()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(bkAddresses))
	for _, v := range bkAddresses {
		go func(a dao.Address) {
			defer wg.Done()

			exists := false
			for _, address := range addresses {
				if a.Id == address.Id {
					exists = true
				}
			}
			if !exists {
				record := dao.Record[dao.Address]{Id: a.Id}
				if err := dao.DeleteRecord(record, dao.DELETE_ADDRESS); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}
	wg.Wait()
}

func deleteRolesFromBackup(roles []dao.Role) {
	bkRoles, err := dao.FindAllRoles()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(bkRoles))
	for _, v := range bkRoles {
		go func(r dao.Role) {
			defer wg.Done()

			exists := false
			for _, role := range roles {
				if r.Id == role.Id {
					exists = true
				}
			}
			if !exists {
				record := dao.Record[dao.Role]{Id: r.Id}
				if err := dao.DeleteRecord(record, dao.DELETE_ROLE); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}
	wg.Wait()
}

func deleteUsersFromBackup(users []dao.User) {
	bkUsers, err := dao.FindAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(bkUsers))

	for _, v := range bkUsers {
		go func(u dao.User) {
			defer wg.Done()

			exists := false
			for _, user := range users {
				if u.Id == user.Id {
					exists = true
				}
			}
			if !exists {
				record := dao.Record[dao.User]{Id: u.Id}
				if err := dao.DeleteRecord(record, dao.DELETE_USER); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}

	wg.Wait()
}

func deleteUserRolesFromBackup(urs []dao.UrXref) {
	bkUrs, err := dao.FindAllXrefs[dao.UrXref](dao.FINDALL_UR)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(bkUrs))
	for _, v := range bkUrs {
		go func(ur dao.XrefRecord[dao.UrXref]) {
			defer wg.Done()

			exists := false
			for _, userRole := range urs {
				if ur.Id == userRole.Id {
					exists = true
				}
			}
			if !exists {
				r := dao.Record[dao.UrXref]{Id: ur.Id}
				if err := dao.DeleteRecord(r, dao.DELETE_UR); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}
	wg.Wait()
}

func deleteUserAddressesFromBackup(uas []dao.UaXref) {
	bkUas, err := dao.FindAllXrefs[dao.UaXref](dao.FINDALL_UA)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(bkUas))

	for _, v := range bkUas {
		go func(ua dao.XrefRecord[dao.UaXref]) {
			defer wg.Done()

			exists := false
			for _, userAddress := range uas {
				if ua.Id == userAddress.Id {
					exists = true
				}
			}
			if !exists {
				r := dao.Record[dao.UaXref]{Id: ua.Id}
				if err := dao.DeleteRecord(r, dao.DELETE_UA); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}
	wg.Wait()
}

func deleteUserPhonesFromBackup(ups []dao.UpXref) {
	bkUps, err := dao.FindAllXrefs[dao.UpXref](dao.FINDALL_UP)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(bkUps))

	for _, v := range bkUps {
		go func(up dao.XrefRecord[dao.UpXref]) {
			defer wg.Done()

			exists := false
			for _, userPhone := range ups {
				if up.Id == userPhone.Id {
					exists = true
				}
			}
			if !exists {
				r := dao.Record[dao.UpXref]{Id: up.Id}
				if err := dao.DeleteRecord(r, dao.DELETE_UP); err != nil {
					log.Fatal(err)
				}
			}
		}(v)
	}
	wg.Wait()
}

func reconstructRoleTables(users []dao.User) (urs []dao.UrXref, rs []dao.Role) {

	for _, v := range users {
		// user_role table records unique to user: dupe check unnecessary
		for _, userRole := range v.UserRoles {
			ur := dao.UrXref{Id: userRole.Id, User_id: v.Id, Role_id: userRole.Role.Id}
			urs = append(urs, ur)
		}

		// dupe check necessary for roles
		for _, ur := range v.UserRoles {
			exists := false
			for _, role := range rs {
				if ur.Role.Id == role.Id {
					exists = true
				}
			}
			if !exists {
				rs = append(rs, ur.Role)
			}
		}
	}
	return urs, rs
}

func reconstructAddressTables(users []dao.User) (uas []dao.UaXref, as []dao.Address) {

	for _, v := range users {

		for _, userAddress := range v.UserAddresses {
			ua := dao.UaXref{Id: userAddress.Id, User_id: v.Id, Address_id: userAddress.Address.Id}
			uas = append(uas, ua)
		}

		for _, ua := range v.UserAddresses {
			exists := false
			for _, address := range as {
				if ua.Address.Id == address.Id {
					exists = true
				}
			}
			if !exists {
				as = append(as, ua.Address)
			}
		}
	}
	return uas, as
}

func reconstructPhoneTables(users []dao.User) (ups []dao.UpXref, ps []dao.Phone) {

	for _, v := range users {

		for _, userPhone := range v.UserPhones {
			up := dao.UpXref{Id: userPhone.Id, User_id: v.Id, Phone_id: userPhone.Phone.Id}
			ups = append(ups, up)
		}

		for _, up := range v.UserPhones {
			exists := false
			for _, phone := range ps {
				if up.Phone.Id == phone.Id {
					exists = true
				}
			}
			if !exists {
				ps = append(ps, up.Phone)
			}
		}
	}
	return ups, ps
}
