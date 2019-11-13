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

	_, err = db.Database.DB().Exec(db.DDL.SecretService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableSecret, err)
	}
}

func TestDatabase_Client_GetSecret_Org(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "*"
	name := "bar"
	value := "baz"
	typee := "org"
	want := testSecret()
	want.ID = &one
	want.Org = &org
	want.Repo = &repo
	want.Name = &name
	want.Value = &value
	want.Type = &typee

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()
	_ = db.CreateSecret(want)

	// run test
	got, err := db.GetSecret(typee, org, repo, name)

	if err != nil {
		t.Errorf("GetSecret returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetSecret is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetSecret_Repo(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "repo"
	want := testSecret()
	want.ID = &one
	want.Org = &org
	want.Repo = &repo
	want.Name = &name
	want.Value = &value
	want.Type = &typee

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()
	_ = db.CreateSecret(want)

	// run test
	got, err := db.GetSecret(typee, org, repo, name)

	if err != nil {
		t.Errorf("GetSecret returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetSecret is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetSecret_Shared(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	team := "bar"
	name := "baz"
	value := "foob"
	typee := "shared"
	want := testSecret()
	want.ID = &one
	want.Org = &org
	want.Team = &team
	want.Name = &name
	want.Value = &value
	want.Type = &typee

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from secrets;")
		db.Database.Close()
	}()
	_ = db.CreateSecret(want)

	// run test
	got, err := db.GetSecret(typee, org, team, name)

	if err != nil {
		t.Errorf("GetSecret returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetSecret is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetSecretList(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "repo"
	sOne := testSecret()
	sOne.ID = &one
	sOne.Org = &org
	sOne.Repo = &repo
	sOne.Name = &name
	sOne.Value = &value
	sOne.Type = &typee
	two := int64(2)
	sTwo := testSecret()
	sTwo.ID = &two
	sTwo.Org = &org
	sTwo.Repo = &repo
	sTwo.Name = &value
	sTwo.Value = &name
	sTwo.Type = &typee
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
	one := int64(1)
	org := "foo"
	repo := "*"
	name := "bar"
	value := "baz"
	typee := "org"
	sOne := testSecret()
	sOne.ID = &one
	sOne.Org = &org
	sOne.Repo = &repo
	sOne.Name = &name
	sOne.Value = &value
	sOne.Type = &typee
	two := int64(2)
	sTwo := testSecret()
	sTwo.ID = &two
	sTwo.Org = &org
	sTwo.Repo = &repo
	sTwo.Name = &value
	sTwo.Value = &name
	sTwo.Type = &typee
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
	got, err := db.GetTypeSecretList(typee, org, repo, 1, 10)

	if err != nil {
		t.Errorf("GetTypeSecretList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetTypeSecretList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretList_Repo(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "repo"
	sOne := testSecret()
	sOne.ID = &one
	sOne.Org = &org
	sOne.Repo = &repo
	sOne.Name = &name
	sOne.Value = &value
	sOne.Type = &typee
	two := int64(2)
	sTwo := testSecret()
	sTwo.ID = &two
	sTwo.Org = &org
	sTwo.Repo = &repo
	sTwo.Name = &value
	sTwo.Value = &name
	sTwo.Type = &typee
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
	got, err := db.GetTypeSecretList(typee, org, repo, 1, 10)

	if err != nil {
		t.Errorf("GetTypeSecretList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetTypeSecretList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretList_Shared(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	team := "bar"
	name := "baz"
	value := "foob"
	typee := "shared"
	sOne := testSecret()
	sOne.ID = &one
	sOne.Org = &org
	sOne.Team = &team
	sOne.Name = &name
	sOne.Value = &value
	sOne.Type = &typee
	two := int64(2)
	sTwo := testSecret()
	sTwo.ID = &two
	sTwo.Org = &org
	sTwo.Team = &team
	sTwo.Name = &value
	sTwo.Value = &name
	sTwo.Type = &typee
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
	got, err := db.GetTypeSecretList(typee, org, team, 1, 10)

	if err != nil {
		t.Errorf("GetTypeSecretList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetTypeSecretList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretCount_Org(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "*"
	name := "bar"
	value := "baz"
	typee := "org"
	sOne := testSecret()
	sOne.ID = &one
	sOne.Org = &org
	sOne.Repo = &repo
	sOne.Name = &name
	sOne.Value = &value
	sOne.Type = &typee
	two := int64(2)
	sTwo := testSecret()
	sTwo.ID = &two
	sTwo.Org = &org
	sTwo.Repo = &repo
	sTwo.Name = &value
	sTwo.Value = &name
	sTwo.Type = &typee
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
	got, err := db.GetTypeSecretCount(typee, org, repo)

	if err != nil {
		t.Errorf("GetTypeSecretCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetTypeSecretCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretCount_Repo(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "repo"
	sOne := testSecret()
	sOne.ID = &one
	sOne.Org = &org
	sOne.Repo = &repo
	sOne.Name = &name
	sOne.Value = &value
	sOne.Type = &typee
	two := int64(2)
	sTwo := testSecret()
	sTwo.ID = &two
	sTwo.Org = &org
	sTwo.Repo = &repo
	sTwo.Name = &value
	sTwo.Value = &name
	sTwo.Type = &typee
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
	got, err := db.GetTypeSecretCount(typee, org, repo)

	if err != nil {
		t.Errorf("GetTypeSecretCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetTypeSecretCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetTypeSecretCount_Shared(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	team := "bar"
	name := "baz"
	value := "foob"
	typee := "shared"
	sOne := testSecret()
	sOne.ID = &one
	sOne.Org = &org
	sOne.Team = &team
	sOne.Name = &name
	sOne.Value = &value
	sOne.Type = &typee
	two := int64(2)
	sTwo := testSecret()
	sTwo.ID = &two
	sTwo.Org = &org
	sTwo.Team = &team
	sTwo.Name = &value
	sTwo.Value = &name
	sTwo.Type = &typee
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
	got, err := db.GetTypeSecretCount(typee, org, team)

	if err != nil {
		t.Errorf("GetTypeSecretCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetTypeSecretCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateSecret(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "repo"
	want := testSecret()
	want.ID = &one
	want.Org = &org
	want.Repo = &repo
	want.Name = &name
	want.Value = &value
	want.Type = &typee

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

	got, _ := db.GetSecret(typee, org, repo, name)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateSecret is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateSecret_Invalid(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	value := "foob"
	typee := "repo"
	s := testSecret()
	s.ID = &one
	s.Org = &org
	s.Repo = &repo
	s.Value = &value
	s.Type = &typee

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
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "repo"
	want := testSecret()
	want.ID = &one
	want.Org = &org
	want.Repo = &repo
	want.Name = &name
	want.Value = &value
	want.Type = &typee

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

	got, _ := db.GetSecret(typee, org, repo, name)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateSecret is %v, want %v", got, want)
	}
}

func TestDatabase_Client_UpdateSecret_Invalid(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	value := "foob"
	typee := "repo"
	s := testSecret()
	s.ID = &one
	s.Org = &org
	s.Repo = &repo
	s.Value = &value
	s.Type = &typee

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
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "repo"
	want := testSecret()
	want.ID = &one
	want.Org = &org
	want.Repo = &repo
	want.Name = &name
	want.Value = &value
	want.Type = &typee

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
