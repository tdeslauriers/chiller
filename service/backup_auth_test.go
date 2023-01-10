package service

import (
	"chiller/http_client"
	"fmt"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestBackupAuthService(t *testing.T) {

	bearer, _ := http_client.GetBearerToken()
	BackupAuthService(bearer)
}

// func TestReconileRoles(t *testing.T) {

// 	auth, err := http_client.GetAuthServiceData()
// 	if err != nil {
// 		panic(err)
// 	}

// 	reconcileRoles(auth)
// }

// func TestReconileUserRoles(t *testing.T) {

// 	auth, err := http_client.GetAuthServiceData()
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, v := range auth {

// 		t.Log(reconcileUserRoles(v))
// 	}
// }

func TestGoPractice(t *testing.T) {

	sl := []string{"atomic", "dog", "bow", "wow", "yippie-yo", "yippie-yay"}

	var wg sync.WaitGroup
	wg.Add(len(sl))

	t.Log("Starting go routines")

	for _, v := range sl {

		go func(f string) {
			defer wg.Done()
			t.Logf("Walking the dog: %s", f)
		}(v)
	}

	t.Log("Fires before wg.Wait")

	wg.Wait()

	t.Log("Loop complete.")
}

func TestGenerics(t *testing.T) {

	ac := animal[cat]{}
	ac.says("meow")

	ad := animal[dog]{}
	ad.says("woof")

}

type talk interface {
	cat | dog
}

type cat struct {
	name string
}

type dog struct {
	name string
}

type animal[T talk] struct {
}

func (a animal[T]) says(sound string) {
	fmt.Print(sound)
}

func TestReconstructRoleTables(t *testing.T) {

	bearer, _ := http_client.GetBearerToken()
	users, err := http_client.GetAuthServiceData(bearer)
	if err != nil {
		panic(err)
	}
	urs, rs := reconstructRoleTables(users)
	t.Log(urs)
	t.Log(rs)

	uas, as := reconstructAddressTables(users)
	t.Log(uas)
	t.Log(as)

	ups, ps := reconstructPhoneTables(users)
	t.Log(ups)
	t.Log(ps)
}
