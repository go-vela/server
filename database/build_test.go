// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"fmt"
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

	_, err = db.Database.DB().Exec(db.DDL.BuildService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableBuild, err)
	}
}

func TestDatabase_Client_GetBuild(t *testing.T) {
	// setup types
	id := int64(1)
	org := "foo"
	repo := "bar"
	name := fmt.Sprintf("%s/%s", org, repo)
	r := testRepo()
	r.ID = &id
	r.Org = &org
	r.Name = &repo
	r.FullName = &name
	oneNum := 1
	want := testBuild()
	want.ID = &id
	want.RepoID = &id
	want.Number = &oneNum

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()
	_ = database.CreateBuild(want)

	// run test
	got, err := database.GetBuild(int(want.GetNumber()), r)

	if err != nil {
		t.Errorf("GetBuild returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBuild is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetLastBuild(t *testing.T) {
	// setup types
	id := int64(1)
	org := "foo"
	repo := "bar"
	name := fmt.Sprintf("%s/%s", org, repo)
	r := testRepo()
	r.ID = &id
	r.Org = &org
	r.Name = &repo
	r.FullName = &name
	oneNum := 1
	b := testBuild()
	b.ID = &id
	b.RepoID = &id
	b.Number = &oneNum
	two := int64(2)
	twoNum := 2
	want := testBuild()
	want.ID = &two
	want.RepoID = &id
	want.Number = &twoNum

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()
	_ = database.CreateBuild(b)
	_ = database.CreateBuild(want)

	// run test
	got, err := database.GetLastBuild(r)

	if err != nil {
		t.Errorf("GetLastBuild returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetLastBuild is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetLastBuild_NotFound(t *testing.T) {
	// setup types
	id := int64(1)
	org := "foo"
	repo := "bar"
	name := fmt.Sprintf("%s/%s", org, repo)
	r := testRepo()
	r.ID = &id
	r.Org = &org
	r.Name = &repo
	r.FullName = &name

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()

	// run test
	got, err := database.GetLastBuild(r)

	if err != nil {
		t.Errorf("GetLastBuild returned err: %v", err)
	}

	if got != nil {
		t.Errorf("GetLastBuild is %v, want nil", got)
	}
}

func TestDatabase_Client_GetBuildList(t *testing.T) {
	// setup types
	one := int64(1)
	oneNum := 1
	bOne := testBuild()
	bOne.ID = &one
	bOne.RepoID = &one
	bOne.Number = &oneNum
	two := int64(2)
	twoNum := 2
	bTwo := testBuild()
	bTwo.ID = &two
	bTwo.RepoID = &one
	bTwo.Number = &twoNum
	want := []*library.Build{bOne, bTwo}

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()
	_ = database.CreateBuild(bOne)
	_ = database.CreateBuild(bTwo)

	// run test
	got, err := database.GetBuildList()

	if err != nil {
		t.Errorf("GetBuildList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBuildList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetBuildCount(t *testing.T) {
	// setup types
	one := int64(1)
	oneNum := 1
	bOne := testBuild()
	bOne.ID = &one
	bOne.RepoID = &one
	bOne.Number = &oneNum
	two := int64(2)
	twoNum := 2
	bTwo := testBuild()
	bTwo.ID = &two
	bTwo.RepoID = &one
	bTwo.Number = &twoNum
	want := 2

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()
	_ = database.CreateBuild(bOne)
	_ = database.CreateBuild(bTwo)

	// run test
	got, err := database.GetBuildCount()

	if err != nil {
		t.Errorf("GetBuildCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetBuildCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetBuildCountByStatus(t *testing.T) {
	// setup types
	one := int64(1)
	oneNum := 1
	pStatus := "pending"
	bOne := testBuild()
	bOne.ID = &one
	bOne.RepoID = &one
	bOne.Number = &oneNum
	bOne.Status = &pStatus
	two := int64(2)
	twoNum := 2
	rStatus := "running"
	bTwo := testBuild()
	bTwo.ID = &two
	bTwo.RepoID = &one
	bTwo.Number = &twoNum
	bTwo.Status = &rStatus
	want := 1

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()
	_ = database.CreateBuild(bOne)
	_ = database.CreateBuild(bTwo)

	// run test
	got, err := database.GetBuildCountByStatus("running")

	if err != nil {
		t.Errorf("GetBuildCountByStatus returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetBuildCountByStatus is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetRepoBuildList(t *testing.T) {
	// setup types
	id := int64(1)
	org := "foo"
	repo := "bar"
	name := fmt.Sprintf("%s/%s", org, repo)
	r := testRepo()
	r.ID = &id
	r.Org = &org
	r.Name = &repo
	r.FullName = &name
	one := int64(1)
	oneNum := 1
	pStatus := "pending"
	bOne := testBuild()
	bOne.ID = &one
	bOne.RepoID = &id
	bOne.Number = &oneNum
	bOne.Status = &pStatus
	two := int64(2)
	twoNum := 2
	rStatus := "running"
	bTwo := testBuild()
	bTwo.ID = &two
	bTwo.RepoID = &id
	bTwo.Number = &twoNum
	bTwo.Status = &rStatus
	want := []*library.Build{bTwo, bOne}

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()
	_ = database.CreateBuild(bOne)
	_ = database.CreateBuild(bTwo)

	// run test
	got, err := database.GetRepoBuildList(r, 1, 10)

	if err != nil {
		t.Errorf("GetRepoBuildList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRepoBuildList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetRepoBuildCount(t *testing.T) {
	// setup types
	id := int64(1)
	org := "foo"
	repo := "bar"
	name := fmt.Sprintf("%s/%s", org, repo)
	r := testRepo()
	r.ID = &id
	r.Org = &org
	r.Name = &repo
	r.FullName = &name
	one := int64(1)
	oneNum := 1
	bOne := testBuild()
	bOne.ID = &one
	bOne.RepoID = &id
	bOne.Number = &oneNum
	two := int64(2)
	twoNum := 2
	bTwo := testBuild()
	bTwo.ID = &two
	bTwo.RepoID = &id
	bTwo.Number = &twoNum
	three := int64(2)
	threeNum := 2
	bThree := testBuild()
	bThree.ID = &three
	bThree.RepoID = &id
	bThree.Number = &threeNum
	want := 2

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()
	_ = database.CreateBuild(bOne)
	_ = database.CreateBuild(bTwo)
	_ = database.CreateBuild(bThree)

	// run test
	got, err := database.GetRepoBuildCount(r)

	if err != nil {
		t.Errorf("GetRepoBuildCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetRepoBuildCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateBuild(t *testing.T) {
	// setup types
	id := int64(1)
	org := "foo"
	repo := "bar"
	name := fmt.Sprintf("%s/%s", org, repo)
	r := testRepo()
	r.ID = &id
	r.Org = &org
	r.Name = &repo
	r.FullName = &name
	oneNum := 1
	want := testBuild()
	want.ID = &id
	want.RepoID = &id
	want.Number = &oneNum

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()

	// run test
	err := database.CreateBuild(want)

	if err != nil {
		t.Errorf("CreateBuild returned err: %v", err)
	}

	got, _ := database.GetBuild(int(want.GetNumber()), r)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateBuild is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateBuild_Invalid(t *testing.T) {
	// setup types
	id := int64(1)
	oneNum := 1
	b := testBuild()
	b.ID = &id
	b.Number = &oneNum

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()

	// run test
	err := database.CreateBuild(b)

	if err == nil {
		t.Errorf("CreateBuild should have returned err")
	}
}

func TestDatabase_Client_UpdateBuild(t *testing.T) {
	// setup types
	id := int64(1)
	org := "foo"
	repo := "bar"
	name := fmt.Sprintf("%s/%s", org, repo)
	r := testRepo()
	r.ID = &id
	r.Org = &org
	r.Name = &repo
	r.FullName = &name
	oneNum := 1
	want := testBuild()
	want.ID = &id
	want.RepoID = &id
	want.Number = &oneNum

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()
	_ = database.CreateBuild(want)

	// run test
	err := database.UpdateBuild(want)

	if err != nil {
		t.Errorf("UpdateBuild returned err: %v", err)
	}

	got, _ := database.GetBuild(int(want.GetNumber()), r)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateBuild is %v, want %v", got, want)
	}
}

func TestDatabase_Client_UpdateBuild_Invalid(t *testing.T) {
	// setup types
	id := int64(1)
	oneNum := 1
	b := testBuild()
	b.ID = &id
	b.Number = &oneNum

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()
	_ = database.CreateBuild(b)

	// run test
	err := database.UpdateBuild(b)

	if err == nil {
		t.Errorf("UpdateBuild should have returned err")
	}
}

func TestDatabase_Client_DeleteBuild(t *testing.T) {
	// setup types
	id := int64(1)
	oneNum := 1
	b := testBuild()
	b.ID = &id
	b.RepoID = &id
	b.Number = &oneNum

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()
	_ = database.CreateBuild(b)

	// run test
	err := database.DeleteBuild(b.GetID())

	if err != nil {
		t.Errorf("DeleteBuild returned err: %v", err)
	}
}

// testBuild is a test helper function to create a
// library Build type with all fields set to their
// zero values.
func testBuild() *library.Build {
	i64 := int64(0)
	i := 0
	str := ""
	return &library.Build{
		ID:           &i64,
		RepoID:       &i64,
		Number:       &i,
		Parent:       &i,
		Event:        &str,
		Status:       &str,
		Error:        &str,
		Enqueued:     &i64,
		Created:      &i64,
		Started:      &i64,
		Finished:     &i64,
		Deploy:       &str,
		Clone:        &str,
		Source:       &str,
		Title:        &str,
		Message:      &str,
		Commit:       &str,
		Sender:       &str,
		Author:       &str,
		Branch:       &str,
		Ref:          &str,
		BaseRef:      &str,
		Host:         &str,
		Runtime:      &str,
		Distribution: &str,
	}
}
