// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

	_, err = db.Database.DB().Exec(db.DDL.SecretService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableSecret, err)
	}
}

func TestDatabase_Client_GetSecret_Org(t *testing.T) {
	// setup types
	want := testSecret()
	want.SetID(1)
	want.SetOrg("foo")
	want.SetRepo("*")
	want.SetName("bar")
	want.SetValue("baz")
	want.SetType("org")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(want)

	// run test
	got, err := db.GetSecret("org", "foo", "*", "bar")

	if err != nil {
		t.Errorf("GetSecret returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetSecret is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetSecret_Repo(t *testing.T) {
	// setup types
	want := testSecret()
	want.SetID(1)
	want.SetOrg("foo")
	want.SetRepo("bar")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("repo")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(want)

	// run test
	got, err := db.GetSecret("repo", "foo", "bar", "baz")

	if err != nil {
		t.Errorf("GetSecret returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetSecret is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetSecret_Shared(t *testing.T) {
	// setup types
	want := testSecret()
	want.SetID(1)
	want.SetOrg("foo")
	want.SetTeam("bar")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("shared")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(want)

	// run test
	got, err := db.GetSecret("shared", "foo", "bar", "baz")

	if err != nil {
		t.Errorf("GetSecret returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetSecret is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetSecretList(t *testing.T) {
	// setup types
	sOne := testSecret()
	sOne.SetID(1)
	sOne.SetOrg("foo")
	sOne.SetRepo("bar")
	sOne.SetName("baz")
	sOne.SetValue("foob")
	sOne.SetType("repo")

	sTwo := testSecret()
	sTwo.SetID(2)
	sTwo.SetOrg("foo")
	sTwo.SetRepo("bar")
	sTwo.SetName("foob")
	sTwo.SetValue("baz")
	sTwo.SetType("repo")

	want := []*library.Secret{sOne, sTwo}

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(sOne)
	_ = db.CreateSecret(sTwo)

	// run test
	got, err := db.GetSecretList()

	if err != nil {
		t.Errorf("GetSecretList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetSecretList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretList_Org(t *testing.T) {
	// setup types
	sOne := testSecret()
	sOne.SetID(1)
	sOne.SetOrg("foo")
	sOne.SetRepo("*")
	sOne.SetName("bar")
	sOne.SetValue("baz")
	sOne.SetType("org")

	sTwo := testSecret()
	sTwo.SetID(2)
	sTwo.SetOrg("foo")
	sTwo.SetRepo("*")
	sTwo.SetName("baz")
	sTwo.SetValue("bar")
	sTwo.SetType("org")

	want := []*library.Secret{sTwo, sOne}

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(sOne)
	_ = db.CreateSecret(sTwo)

	// run test
	got, err := db.GetTypeSecretList("org", "foo", "*", 1, 10)

	if err != nil {
		t.Errorf("GetTypeSecretList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetTypeSecretList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretList_Repo(t *testing.T) {
	// setup types
	sOne := testSecret()
	sOne.SetID(1)
	sOne.SetOrg("foo")
	sOne.SetRepo("bar")
	sOne.SetName("baz")
	sOne.SetValue("foob")
	sOne.SetType("repo")

	sTwo := testSecret()
	sTwo.SetID(2)
	sTwo.SetOrg("foo")
	sTwo.SetRepo("bar")
	sTwo.SetName("foob")
	sTwo.SetValue("baz")
	sTwo.SetType("repo")

	want := []*library.Secret{sTwo, sOne}

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(sOne)
	_ = db.CreateSecret(sTwo)

	// run test
	got, err := db.GetTypeSecretList("repo", "foo", "bar", 1, 10)

	if err != nil {
		t.Errorf("GetTypeSecretList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetTypeSecretList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretList_Shared(t *testing.T) {
	// setup types
	sOne := testSecret()
	sOne.SetID(1)
	sOne.SetOrg("foo")
	sOne.SetTeam("bar")
	sOne.SetName("baz")
	sOne.SetValue("foob")
	sOne.SetType("shared")

	sTwo := testSecret()
	sTwo.SetID(2)
	sTwo.SetOrg("foo")
	sTwo.SetTeam("bar")
	sTwo.SetName("foob")
	sTwo.SetValue("baz")
	sTwo.SetType("shared")

	want := []*library.Secret{sTwo, sOne}

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(sOne)
	_ = db.CreateSecret(sTwo)

	// run test
	got, err := db.GetTypeSecretList("shared", "foo", "bar", 1, 10)

	if err != nil {
		t.Errorf("GetTypeSecretList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetTypeSecretList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretCount_Org(t *testing.T) {
	// setup types
	sOne := testSecret()
	sOne.SetID(1)
	sOne.SetOrg("foo")
	sOne.SetRepo("*")
	sOne.SetName("bar")
	sOne.SetValue("baz")
	sOne.SetType("org")

	sTwo := testSecret()
	sTwo.SetID(2)
	sTwo.SetOrg("foo")
	sTwo.SetRepo("*")
	sTwo.SetName("baz")
	sTwo.SetValue("bar")
	sTwo.SetType("org")

	want := 2

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(sOne)
	_ = db.CreateSecret(sTwo)

	// run test
	got, err := db.GetTypeSecretCount("org", "foo", "*")

	if err != nil {
		t.Errorf("GetTypeSecretCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetTypeSecretCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretCount_Repo(t *testing.T) {
	// setup types
	sOne := testSecret()
	sOne.SetID(1)
	sOne.SetOrg("foo")
	sOne.SetRepo("bar")
	sOne.SetName("baz")
	sOne.SetValue("foob")
	sOne.SetType("repo")

	sTwo := testSecret()
	sTwo.SetID(2)
	sTwo.SetOrg("foo")
	sTwo.SetRepo("bar")
	sTwo.SetName("foob")
	sTwo.SetValue("baz")
	sTwo.SetType("repo")

	want := 2

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(sOne)
	_ = db.CreateSecret(sTwo)

	// run test
	got, err := db.GetTypeSecretCount("repo", "foo", "bar")

	if err != nil {
		t.Errorf("GetTypeSecretCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetTypeSecretCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretCount_Shared(t *testing.T) {
	// setup types
	sOne := testSecret()
	sOne.SetID(1)
	sOne.SetOrg("foo")
	sOne.SetTeam("bar")
	sOne.SetName("baz")
	sOne.SetValue("foob")
	sOne.SetType("shared")

	sTwo := testSecret()
	sTwo.SetID(2)
	sTwo.SetOrg("foo")
	sTwo.SetTeam("bar")
	sTwo.SetName("foob")
	sTwo.SetValue("baz")
	sTwo.SetType("shared")

	want := 2

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(sOne)
	_ = db.CreateSecret(sTwo)

	// run test
	got, err := db.GetTypeSecretCount("shared", "foo", "bar")

	if err != nil {
		t.Errorf("GetTypeSecretCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetTypeSecretCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateSecret(t *testing.T) {
	// setup types
	want := testSecret()
	want.SetID(1)
	want.SetOrg("foo")
	want.SetRepo("bar")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("repo")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateSecret(want)

	if err != nil {
		t.Errorf("CreateSecret returned err: %v", err)
	}

	got, _ := db.GetSecret("repo", "foo", "bar", "baz")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateSecret is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateSecret_Invalid(t *testing.T) {
	// setup types
	s := testSecret()
	s.SetID(1)
	s.SetOrg("foo")
	s.SetRepo("bar")
	s.SetValue("foob")
	s.SetType("repo")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateSecret(s)

	if err == nil {
		t.Errorf("CreateSecret should have returned err")
	}
}

func TestDatabase_Client_UpdateSecret(t *testing.T) {
	// setup types
	want := testSecret()
	want.SetID(1)
	want.SetOrg("foo")
	want.SetRepo("bar")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("repo")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(want)

	// run test
	err := db.UpdateSecret(want)

	if err != nil {
		t.Errorf("UpdateSecret returned err: %v", err)
	}

	got, _ := db.GetSecret("repo", "foo", "bar", "baz")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateSecret is %v, want %v", got, want)
	}
}

func TestDatabase_Client_UpdateSecret_Invalid(t *testing.T) {
	// setup types
	s := testSecret()
	s.SetID(1)
	s.SetOrg("foo")
	s.SetRepo("bar")
	s.SetValue("foob")
	s.SetType("repo")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(s)

	// run test
	err := db.UpdateSecret(s)

	if err == nil {
		t.Errorf("UpdateSecret should have returned err")
	}
}

func TestDatabase_Client_DeleteSecret(t *testing.T) {
	// setup types
	want := testSecret()
	want.SetID(1)
	want.SetOrg("foo")
	want.SetRepo("bar")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("repo")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()

	_ = db.CreateSecret(want)

	// run test
	err := db.DeleteSecret(want.GetID())

	if err != nil {
		t.Errorf("DeleteSecret returned err: %v", err)
	}
}

// testSecret is a test helper function to create a
// library Secret type with all fields set to their
// zero values.
func testSecret() *library.Secret {
	i64 := int64(0)
	str := ""
	arr := []string{}
	booL := false

	return &library.Secret{
		ID:           &i64,
		Org:          &str,
		Repo:         &str,
		Team:         &str,
		Name:         &str,
		Value:        &str,
		Type:         &str,
		Images:       &arr,
		Events:       &arr,
		AllowCommand: &booL,
	}
}
