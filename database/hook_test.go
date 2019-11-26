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

	_, err = db.Database.DB().Exec(db.DDL.HookService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableHook, err)
	}
}

func TestDatabase_Client_GetHook(t *testing.T) {
	// setup types
	r := testRepo()
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	want := testHook()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from hooks;")
		db.Database.Close()
	}()
	_ = db.CreateHook(want)

	// run test
	got, err := db.GetHook(want.GetNumber(), r)

	if err != nil {
		t.Errorf("GetHook returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHook is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetHookList(t *testing.T) {
	// setup types
	hOne := testHook()
	hOne.SetID(1)
	hOne.SetRepoID(1)
	hOne.SetBuildID(1)
	hOne.SetNumber(1)
	hOne.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	hTwo := testHook()
	hTwo.SetID(2)
	hTwo.SetRepoID(1)
	hTwo.SetBuildID(2)
	hTwo.SetNumber(2)
	hTwo.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	want := []*library.Hook{hOne, hTwo}

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from hooks;")
		db.Database.Close()
	}()
	_ = db.CreateHook(hOne)
	_ = db.CreateHook(hTwo)

	// run test
	got, err := db.GetHookList()

	if err != nil {
		t.Errorf("GetHookList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetHookList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetRepoHookList(t *testing.T) {
	// setup types
	hOne := testHook()
	hOne.SetID(1)
	hOne.SetRepoID(1)
	hOne.SetBuildID(1)
	hOne.SetNumber(1)
	hOne.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	hTwo := testHook()
	hTwo.SetID(2)
	hTwo.SetRepoID(1)
	hTwo.SetBuildID(2)
	hTwo.SetNumber(2)
	hTwo.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	r := testRepo()
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	want := []*library.Hook{hTwo, hOne}

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from hooks;")
		db.Database.Close()
	}()
	_ = db.CreateHook(hOne)
	_ = db.CreateHook(hTwo)

	// run test
	got, err := db.GetRepoHookList(r, 1, 10)

	if err != nil {
		t.Errorf("GetRepoHookList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRepoHookList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetRepoHookCount(t *testing.T) {
	// setup types
	hOne := testHook()
	hOne.SetID(1)
	hOne.SetRepoID(1)
	hOne.SetBuildID(1)
	hOne.SetNumber(1)
	hOne.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	hTwo := testHook()
	hTwo.SetID(2)
	hTwo.SetRepoID(1)
	hTwo.SetBuildID(2)
	hTwo.SetNumber(2)
	hTwo.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	r := testRepo()
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	want := 2

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from hooks;")
		db.Database.Close()
	}()
	_ = db.CreateHook(hOne)
	_ = db.CreateHook(hTwo)

	// run test
	got, err := db.GetRepoHookCount(r)

	if err != nil {
		t.Errorf("GetRepoHookCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetRepoHookCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateHook(t *testing.T) {
	// setup types
	want := testHook()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from hooks;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateHook(want)

	if err != nil {
		t.Errorf("CreateHook returned err: %v", err)
	}

	got, _ := db.GetHookList()

	if !reflect.DeepEqual(got[0], want) {
		t.Errorf("CreateHook is %v, want %v", got[0], want)
	}
}

func TestDatabase_Client_CreateHook_Invalid(t *testing.T) {
	// setup types
	h := testHook()
	h.SetID(1)
	h.SetBuildID(1)
	h.SetNumber(1)
	h.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from hooks;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateHook(h)

	if err == nil {
		t.Errorf("CreateHook should have returned err")
	}
}

func TestDatabase_Client_UpdateHook(t *testing.T) {
	// setup types
	want := testHook()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from hooks;")
		db.Database.Close()
	}()
	_ = db.CreateHook(want)

	// run test
	err := db.UpdateHook(want)

	if err != nil {
		t.Errorf("UpdateHook returned err: %v", err)
	}

	got, _ := db.GetHookList()

	if !reflect.DeepEqual(got[0], want) {
		t.Errorf("UpdateHook is %v, want %v", got[0], want)
	}
}

func TestDatabase_Client_UpdateHook_Invalid(t *testing.T) {
	// setup types
	h := testHook()
	h.SetID(1)
	h.SetBuildID(1)
	h.SetNumber(1)
	h.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from hooks;")
		db.Database.Close()
	}()
	_ = db.CreateHook(h)

	// run test
	err := db.UpdateHook(h)

	if err == nil {
		t.Errorf("UpdateHook should have returned err")
	}

}

func TestDatabase_Client_DeleteHook(t *testing.T) {
	// setup types
	h := testHook()
	h.SetID(1)
	h.SetRepoID(1)
	h.SetBuildID(1)
	h.SetNumber(1)
	h.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from hooks;")
		db.Database.Close()
	}()
	_ = db.CreateHook(h)

	// run test
	err := db.DeleteHook(h.GetID())

	if err != nil {
		t.Errorf("DeleteHook returned err: %v", err)
	}
}

// testHook is a test helper function to create a
// library Hook type with all fields set to their
// zero values.
func testHook() *library.Hook {
	i := 0
	i64 := int64(0)
	str := ""
	return &library.Hook{
		ID:       &i64,
		RepoID:   &i64,
		BuildID:  &i64,
		Number:   &i,
		SourceID: &str,
		Created:  &i64,
		Host:     &str,
		Event:    &str,
		Branch:   &str,
		Error:    &str,
		Status:   &str,
		Link:     &str,
	}
}
