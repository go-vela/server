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

	_, err = db.Database.DB().Exec(db.DDL.BuildService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableBuild, err)
	}
}

func TestDatabase_Client_GetBuild(t *testing.T) {
	// setup types
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	want := testBuild()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetNumber(1)

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
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	b := testBuild()
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	want := testBuild()
	want.SetID(2)
	want.SetRepoID(1)
	want.SetNumber(2)

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
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

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
	bOne := testBuild()
	bOne.SetID(1)
	bOne.SetRepoID(1)
	bOne.SetNumber(1)

	bTwo := testBuild()
	bTwo.SetID(2)
	bTwo.SetRepoID(1)
	bTwo.SetNumber(2)

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
	bOne := testBuild()
	bOne.SetID(1)
	bOne.SetRepoID(1)
	bOne.SetNumber(1)

	bTwo := testBuild()
	bTwo.SetID(2)
	bTwo.SetRepoID(1)
	bTwo.SetNumber(2)

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
	bOne := testBuild()
	bOne.SetID(1)
	bOne.SetRepoID(1)
	bOne.SetNumber(1)
	bOne.SetStatus(constants.StatusPending)

	bTwo := testBuild()
	bTwo.SetID(2)
	bTwo.SetRepoID(1)
	bTwo.SetNumber(2)
	bTwo.SetStatus(constants.StatusRunning)

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
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	bOne := testBuild()
	bOne.SetID(1)
	bOne.SetRepoID(1)
	bOne.SetNumber(1)
	bOne.SetStatus(constants.StatusPending)

	bTwo := testBuild()
	bTwo.SetID(2)
	bTwo.SetRepoID(1)
	bTwo.SetNumber(2)
	bTwo.SetStatus(constants.StatusRunning)

	want := []*library.Build{bTwo}
	wantCount := int64(2)

	// setup database
	database, _ := NewTest()

	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()

	_ = database.CreateBuild(bOne)
	_ = database.CreateBuild(bTwo)

	// run test
	got, gotCount, err := database.GetRepoBuildList(r, 1, 1)

	if err != nil {
		t.Errorf("GetRepoBuildList returned err: %v", err)
	}

	if gotCount != wantCount {
		t.Errorf("Count for GetRepoBuildList returned %v, want %v", gotCount, wantCount)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRepoBuildList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetRepoBuildListByEvent(t *testing.T) {
	// setup types
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	bOne := testBuild()
	bOne.SetID(1)
	bOne.SetRepoID(1)
	bOne.SetNumber(1)
	bOne.SetEvent("push")
	bOne.SetStatus(constants.StatusPending)

	bTwo := testBuild()
	bTwo.SetID(2)
	bTwo.SetRepoID(1)
	bTwo.SetNumber(2)
	bTwo.SetEvent("tag")
	bTwo.SetStatus(constants.StatusRunning)

	bThree := testBuild()
	bThree.SetID(3)
	bThree.SetRepoID(1)
	bThree.SetNumber(3)
	bThree.SetEvent("push")
	bThree.SetStatus(constants.StatusPending)

	want := []*library.Build{bThree}
	wantCount := int64(2)

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
	got, gotCount, err := database.GetRepoBuildListByEvent(r, 1, 1, "push")

	if err != nil {
		t.Errorf("GetRepoBuildListByEvent returned err: %v", err)
	}

	if gotCount != wantCount {
		t.Errorf("Count for GetRepoBuildListByEvent returned %v, want %v", gotCount, wantCount)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRepoBuildListByEvent is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetRepoBuildListByEvent_No_Results(t *testing.T) {
	// setup types
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	bOne := testBuild()
	bOne.SetID(1)
	bOne.SetRepoID(1)
	bOne.SetNumber(1)
	bOne.SetEvent("push")
	bOne.SetStatus(constants.StatusPending)

	bTwo := testBuild()
	bTwo.SetID(2)
	bTwo.SetRepoID(1)
	bTwo.SetNumber(2)
	bTwo.SetEvent("push")
	bTwo.SetStatus(constants.StatusRunning)

	want := []*library.Build{}
	wantCount := int64(0)

	// setup database
	database, _ := NewTest()

	defer func() {
		database.Database.Exec("delete from builds;")
		database.Database.Close()
	}()

	_ = database.CreateBuild(bOne)
	_ = database.CreateBuild(bTwo)

	// run test
	got, gotCount, err := database.GetRepoBuildListByEvent(r, 1, 1, "tag")

	if err != nil {
		t.Errorf("GetRepoBuildListByEvent returned err: %v", err)
	}

	if gotCount != wantCount {
		t.Errorf("Count for GetRepoBuildListByEvent returned %v, want %v", gotCount, wantCount)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRepoBuildListByEvent is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetRepoBuildCount(t *testing.T) {
	// setup types
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	bOne := testBuild()
	bOne.SetID(1)
	bOne.SetRepoID(1)
	bOne.SetNumber(1)

	bTwo := testBuild()
	bTwo.SetID(2)
	bTwo.SetRepoID(1)
	bTwo.SetNumber(2)

	bThree := testBuild()
	bThree.SetID(3)
	bThree.SetRepoID(2)
	bThree.SetNumber(3)

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

func TestDatabase_Client_GetRepoBuildCountByEvent(t *testing.T) {
	// setup types
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	bOne := testBuild()
	bOne.SetID(1)
	bOne.SetRepoID(1)
	bOne.SetNumber(1)
	bOne.SetEvent("push")

	bTwo := testBuild()
	bTwo.SetID(2)
	bTwo.SetRepoID(1)
	bTwo.SetNumber(2)
	bTwo.SetEvent("tag")

	bThree := testBuild()
	bThree.SetID(3)
	bThree.SetRepoID(1)
	bThree.SetNumber(3)
	bThree.SetEvent("push")

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
	got, err := database.GetRepoBuildCountByEvent(r, "push")

	if err != nil {
		t.Errorf("GetRepoBuildCountByEvent returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetRepoBuildCountByEvent is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateBuild(t *testing.T) {
	// setup types
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	want := testBuild()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetNumber(1)

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
	b := testBuild()
	b.SetID(1)
	b.SetNumber(1)

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
	r := testRepo()
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	want := testBuild()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetNumber(1)

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
	b := testBuild()
	b.SetID(1)
	b.SetNumber(1)

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
	b := testBuild()
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

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
		Email:        &str,
		Link:         &str,
		Branch:       &str,
		Ref:          &str,
		BaseRef:      &str,
		Host:         &str,
		Runtime:      &str,
		Distribution: &str,
	}
}
