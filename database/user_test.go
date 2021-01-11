// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
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

	_, err = db.Database.DB().Exec(db.DDL.UserService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableUser, err)
	}
}

func TestDatabase_Client_GetUser(t *testing.T) {
	// setup types
	want := testUser()
	want.SetID(1)
	want.SetName("foo")
	want.SetToken("bar")
	want.SetHash("baz")
	want.SetFavorites([]string{"foo", "bar"})

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(want)

	// run test
	got, err := db.GetUser(want.GetID())

	if err != nil {
		t.Errorf("GetUser returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetUser is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetUserName(t *testing.T) {
	// setup types
	want := testUser()
	want.SetID(1)
	want.SetName("foo")
	want.SetToken("bar")
	want.SetHash("baz")
	want.SetFavorites([]string{"foo", "bar"})

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(want)

	// run test
	got, err := db.GetUserName(want.GetName())

	if err != nil {
		t.Errorf("GetUserName returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetUserName is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetUserRefreshToken(t *testing.T) {
	// setup types
	want := testUser()
	want.SetID(1)
	want.SetName("foo")
	want.SetToken("bar")
	want.SetRefreshToken("abc")
	want.SetHash("baz")
	want.SetFavorites([]string{"foo", "bar"})

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(want)

	// run test
	got, err := db.GetUserRefreshToken(want.GetRefreshToken())

	if err != nil {
		t.Errorf("GetUserRefreshToken returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetUserRefreshToken is %v, want %v", got, want)
	}
}
func TestDatabase_Client_GetUserCount(t *testing.T) {
	// setup types
	uOne := testUser()
	uOne.SetID(1)
	uOne.SetName("foo")
	uOne.SetToken("bar")
	uOne.SetHash("baz")
	uOne.SetFavorites([]string{"foo", "bar"})

	uTwo := testUser()
	uTwo.SetID(2)
	uTwo.SetName("bar")
	uTwo.SetToken("foo")
	uTwo.SetHash("baz")
	uTwo.SetFavorites([]string{"baz"})

	want := 2

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(uOne)
	_ = db.CreateUser(uTwo)

	// run test
	got, err := db.GetUserCount()

	if err != nil {
		t.Errorf("GetUserCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetUserCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetUserList(t *testing.T) {
	// setup types
	uOne := testUser()
	uOne.SetID(1)
	uOne.SetName("foo")
	uOne.SetToken("bar")
	uOne.SetHash("baz")
	uOne.SetFavorites([]string{"foo", "bar"})

	uTwo := testUser()
	uTwo.SetID(2)
	uTwo.SetName("bar")
	uTwo.SetToken("foo")
	uTwo.SetHash("baz")
	uTwo.SetFavorites([]string{"baz"})

	want := []*library.User{uOne, uTwo}

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(uOne)
	_ = db.CreateUser(uTwo)

	// run test
	got, err := db.GetUserList()

	if err != nil {
		t.Errorf("GetUserList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetUserList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetUserLiteList(t *testing.T) {
	// setup types
	uOne := testUser()
	uOne.SetID(1)
	uOne.SetName("foo")
	uOne.SetToken("bar")
	uOne.SetHash("baz")
	uOne.SetFavorites([]string{"foo"})

	uTwo := testUser()
	uTwo.SetID(2)
	uTwo.SetName("bar")
	uTwo.SetToken("foo")
	uTwo.SetHash("baz")
	uTwo.SetFavorites([]string{"baz"})

	wOne := testUser()
	wOne.SetID(1)
	wOne.SetName("foo")
	wOne.SetFavorites(nil)

	wTwo := testUser()
	wTwo.SetID(2)
	wTwo.SetName("bar")
	wTwo.SetFavorites(nil)

	want := []*library.User{wTwo, wOne}

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(uOne)
	_ = db.CreateUser(uTwo)

	// run test
	got, err := db.GetUserLiteList(1, 10)

	if err != nil {
		t.Errorf("GetUserLiteList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetUserLiteList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateUser(t *testing.T) {
	// setup types
	want := testUser()
	want.SetID(1)
	want.SetName("foo")
	want.SetToken("bar")
	want.SetHash("baz")
	want.SetFavorites([]string{"foo"})

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateUser(want)

	if err != nil {
		t.Errorf("CreateUser returned err: %v", err)
	}

	got, _ := db.GetUser(want.GetID())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateUser is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateUser_Invalid(t *testing.T) {
	// setup types
	u := testUser()
	u.SetID(1)
	u.SetToken("bar")
	u.SetHash("baz")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateUser(u)

	if err == nil {
		t.Errorf("CreateUser should have returned err")
	}
}

func TestDatabase_Client_UpdateUser(t *testing.T) {
	// setup types
	want := testUser()
	want.SetID(1)
	want.SetName("foo")
	want.SetToken("bar")
	want.SetHash("baz")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(want)

	// update favorites
	want.SetFavorites([]string{"foo"})

	// run test
	err := db.UpdateUser(want)

	if err != nil {
		t.Errorf("UpdateUser returned err: %v", err)
	}

	got, _ := db.GetUser(want.GetID())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateUser is %v, want %v", got, want)
	}
}

func TestDatabase_Client_UpdateUser_Invalid(t *testing.T) {
	// setup types
	u := testUser()
	u.SetID(1)
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetFavorites([]string{"foo"})

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(u)

	// run test
	err := db.UpdateUser(u)

	if err == nil {
		t.Errorf("UpdateUser should have returned err")
	}
}

func TestDatabase_Client_DeleteUser(t *testing.T) {
	// setup types
	want := testUser()
	want.SetID(1)
	want.SetName("foo")
	want.SetToken("bar")
	want.SetHash("baz")
	want.SetFavorites([]string{"foo"})

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(want)

	// run test
	err := db.DeleteUser(want.GetID())

	if err != nil {
		t.Errorf("DeleteUser returned err: %v", err)
	}
}

// testUser is a test helper function to create a
// library User type with all fields set to their
// zero values.
func testUser() *library.User {
	i64 := int64(0)
	str := ""
	b := false
	arr := []string{}

	return &library.User{
		ID:           &i64,
		Name:         &str,
		RefreshToken: &str,
		Token:        &str,
		Hash:         &str,
		Favorites:    &arr,
		Active:       &b,
		Admin:        &b,
	}
}
