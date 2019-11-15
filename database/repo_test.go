// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"log"
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

func init() {
	db, err := NewTest()
	if err != nil {
		log.Fatalf("Error creating test database: %v", err)
	}

	_, err = db.Database.DB().Exec(db.DDL.RepoService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableRepo, err)
	}
}

func TestDatabase_Client_GetRepo(t *testing.T) {
	// setup types
	want := testRepo()
	want.SetID(1)
	want.SetUserID(1)
	want.SetOrg("foo")
	want.SetName("bar")
	want.SetFullName("foo/bar")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(want)

	// run test
	got, err := db.GetRepo(want.GetOrg(), want.GetName())

	if err != nil {
		t.Errorf("GetRepo returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRepo is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetRepoList(t *testing.T) {
	// setup types
	rOne := testRepo()
	rOne.SetID(1)
	rOne.SetUserID(1)
	rOne.SetOrg("foo")
	rOne.SetName("bar")
	rOne.SetFullName("foo/bar")

	rTwo := testRepo()
	rTwo.SetID(2)
	rTwo.SetUserID(1)
	rTwo.SetOrg("bar")
	rTwo.SetName("foo")
	rTwo.SetFullName("bar/foo")

	want := []*library.Repo{rOne, rTwo}

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(rOne)
	_ = db.CreateRepo(rTwo)

	// run test
	got, err := db.GetRepoList()

	if err != nil {
		t.Errorf("GetRepoList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRepoList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetRepoCount(t *testing.T) {
	// setup types
	rOne := testRepo()
	rOne.SetID(1)
	rOne.SetUserID(1)
	rOne.SetOrg("foo")
	rOne.SetName("bar")
	rOne.SetFullName("foo/bar")

	rTwo := testRepo()
	rTwo.SetID(2)
	rTwo.SetUserID(1)
	rTwo.SetOrg("bar")
	rTwo.SetName("foo")
	rTwo.SetFullName("bar/foo")

	want := 2

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(rOne)
	_ = db.CreateRepo(rTwo)

	// run test
	got, err := db.GetRepoCount()

	if err != nil {
		t.Errorf("GetRepoCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetRepoCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetUserRepoList(t *testing.T) {
	// setup types
	rOne := testRepo()
	rOne.SetID(1)
	rOne.SetUserID(1)
	rOne.SetOrg("foo")
	rOne.SetName("bar")
	rOne.SetFullName("foo/bar")

	rTwo := testRepo()
	rTwo.SetID(2)
	rTwo.SetUserID(1)
	rTwo.SetOrg("bar")
	rTwo.SetName("foo")
	rTwo.SetFullName("bar/foo")

	u := testUser()
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")

	want := []*library.Repo{rTwo, rOne}

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(rOne)
	_ = db.CreateRepo(rTwo)

	// run test
	got, err := db.GetUserRepoList(u, 1, 10)

	if err != nil {
		t.Errorf("GetUserRepoList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetUserRepoList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetUserRepoCount(t *testing.T) {
	// setup types
	rOne := testRepo()
	rOne.SetID(1)
	rOne.SetUserID(1)
	rOne.SetOrg("foo")
	rOne.SetName("bar")
	rOne.SetFullName("foo/bar")

	rTwo := testRepo()
	rTwo.SetID(2)
	rTwo.SetUserID(1)
	rTwo.SetOrg("bar")
	rTwo.SetName("foo")
	rTwo.SetFullName("bar/foo")

	u := testUser()
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")

	want := 2

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(rOne)
	_ = db.CreateRepo(rTwo)

	// run test
	got, err := db.GetUserRepoCount(u)

	if err != nil {
		t.Errorf("GetUserRepoCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetUserRepoCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateRepo(t *testing.T) {
	// setup types
	want := testRepo()
	want.SetID(1)
	want.SetUserID(1)
	want.SetOrg("foo")
	want.SetName("bar")
	want.SetFullName("foo/bar")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateRepo(want)

	if err != nil {
		t.Errorf("CreateRepo returned err: %v", err)
	}

	got, _ := db.GetRepo(want.GetOrg(), want.GetName())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateRepo is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateRepo_Invalid(t *testing.T) {
	// setup types
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateRepo(r)

	if err == nil {
		t.Errorf("CreateRepo should have returned err")
	}
}

func TestDatabase_Client_UpdateRepo(t *testing.T) {
	// setup types
	want := testRepo()
	want.SetID(1)
	want.SetUserID(1)
	want.SetOrg("foo")
	want.SetName("bar")
	want.SetFullName("foo/bar")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(want)

	// run test
	err := db.UpdateRepo(want)

	if err != nil {
		t.Errorf("UpdateRepo returned err: %v", err)
	}

	got, _ := db.GetRepo(want.GetOrg(), want.GetName())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateRepo is %v, want %v", got, want)
	}
}

func TestDatabase_Client_UpdateRepo_Invalid(t *testing.T) {
	// setup types
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)

	// run test
	err := db.UpdateRepo(r)

	if err == nil {
		t.Errorf("UpdateRepo should have returned err")
	}

}

func TestDatabase_Client_UpdateRepo_Boolean(t *testing.T) {
	// setup types
	want := testRepo()
	want.SetID(1)
	want.SetUserID(1)
	want.SetOrg("foo")
	want.SetName("bar")
	want.SetFullName("foo/bar")
	want.SetAllowPull(true)
	want.SetActive(false)

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(want)

	// run test
	want.SetAllowPull(false)
	want.SetActive(true)

	err := db.UpdateRepo(want)

	if err != nil {
		t.Errorf("UpdateRepo returned err: %+v", err)
	}

	got, _ := db.GetRepo(want.GetOrg(), want.GetName())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateRepo is %+v, want %+v", got, want)
	}
}

func TestDatabase_Client_DeleteRepo(t *testing.T) {
	// setup types
	want := testRepo()
	want.SetID(1)
	want.SetUserID(1)
	want.SetOrg("foo")
	want.SetName("bar")
	want.SetFullName("foo/bar")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(want)

	// run test
	err := db.DeleteRepo(want.GetID())

	if err != nil {
		t.Errorf("DeleteRepo returned err: %v", err)
	}
}

// testRepo is a test helper function to create a
// library Repo type with all fields set to their
// zero values.
func testRepo() *library.Repo {
	i64 := int64(0)
	str := ""
	b := false
	return &library.Repo{
		ID:          &i64,
		UserID:      &i64,
		Org:         &str,
		Name:        &str,
		FullName:    &str,
		Link:        &str,
		Clone:       &str,
		Branch:      &str,
		Timeout:     &i64,
		Visibility:  &str,
		Private:     &b,
		Trusted:     &b,
		Active:      &b,
		AllowPull:   &b,
		AllowPush:   &b,
		AllowDeploy: &b,
		AllowTag:    &b,
	}
}
