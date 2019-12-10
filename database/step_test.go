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

	_, err = db.Database.DB().Exec(db.DDL.StepService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableStep, err)
	}
}

func TestDatabase_Client_GetStep(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	want := testStep()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetName("foo")
	want.SetImage("baz")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()
	_ = db.CreateStep(want)

	// run test
	got, err := db.GetStep(int(want.GetNumber()), b)

	if err != nil {
		t.Errorf("GetStep returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetStep is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetStepList(t *testing.T) {
	// setup types
	sOne := testStep()
	sOne.SetID(1)
	sOne.SetRepoID(1)
	sOne.SetBuildID(1)
	sOne.SetNumber(1)
	sOne.SetName("foo")
	sOne.SetImage("baz")

	sTwo := testStep()
	sTwo.SetID(2)
	sTwo.SetRepoID(1)
	sTwo.SetBuildID(1)
	sTwo.SetNumber(2)
	sTwo.SetName("bar")
	sTwo.SetImage("baz")

	want := []*library.Step{sOne, sTwo}

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()
	_ = db.CreateStep(sOne)
	_ = db.CreateStep(sTwo)

	// run test
	got, err := db.GetStepList()

	if err != nil {
		t.Errorf("GetStepList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetStepList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetBuildStepList(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	sOne := testStep()
	sOne.SetID(1)
	sOne.SetRepoID(1)
	sOne.SetBuildID(1)
	sOne.SetNumber(1)
	sOne.SetName("foo")
	sOne.SetImage("baz")

	sTwo := testStep()
	sTwo.SetID(2)
	sTwo.SetRepoID(1)
	sTwo.SetBuildID(1)
	sTwo.SetNumber(2)
	sTwo.SetName("bar")
	sTwo.SetImage("baz")

	want := []*library.Step{sTwo, sOne}

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()
	_ = db.CreateStep(sOne)
	_ = db.CreateStep(sTwo)

	// run test
	got, err := db.GetBuildStepList(b, 1, 10)

	if err != nil {
		t.Errorf("GetBuildStepList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBuildStepList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetBuildStepCount(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(2)
	b.SetRepoID(1)
	b.SetNumber(1)

	sOne := testStep()
	sOne.SetID(1)
	sOne.SetRepoID(1)
	sOne.SetBuildID(1)
	sOne.SetNumber(1)
	sOne.SetName("foo")
	sOne.SetImage("baz")

	sTwo := testStep()
	sTwo.SetID(2)
	sTwo.SetRepoID(2)
	sTwo.SetBuildID(2)
	sTwo.SetNumber(1)
	sTwo.SetName("foo")
	sTwo.SetImage("baz")

	sThree := testStep()
	sThree.SetID(3)
	sThree.SetRepoID(2)
	sThree.SetBuildID(2)
	sThree.SetNumber(2)
	sThree.SetName("bar")
	sThree.SetImage("baz")

	want := 2

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()
	_ = db.CreateStep(sOne)
	_ = db.CreateStep(sTwo)
	_ = db.CreateStep(sThree)

	// run test
	got, err := db.GetBuildStepCount(b)

	if err != nil {
		t.Errorf("GetBuildStepCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetBuildStepCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetStepImageCount(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(2)
	b.SetRepoID(1)
	b.SetNumber(1)

	sOne := testStep()
	sOne.SetID(1)
	sOne.SetRepoID(1)
	sOne.SetBuildID(1)
	sOne.SetNumber(1)
	sOne.SetName("foo")
	sOne.SetImage("baz")

	sTwo := testStep()
	sTwo.SetID(2)
	sTwo.SetRepoID(2)
	sTwo.SetBuildID(2)
	sTwo.SetNumber(1)
	sTwo.SetName("foo")
	sTwo.SetImage("baz")

	sThree := testStep()
	sThree.SetID(3)
	sThree.SetRepoID(2)
	sThree.SetBuildID(2)
	sThree.SetNumber(2)
	sThree.SetName("bar")
	sThree.SetImage("bazian:latest")

	want := make(map[string]float64)
	want["baz"] = 2
	want["bazian:latest"] = 1

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()
	_ = db.CreateStep(sOne)
	_ = db.CreateStep(sTwo)
	_ = db.CreateStep(sThree)

	// run test
	got, err := db.GetStepImageCount()

	if err != nil {
		t.Errorf("GetStepImageCount returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetStepImageCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetStepStatusCount(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(2)
	b.SetRepoID(1)
	b.SetNumber(1)

	sOne := testStep()
	sOne.SetID(1)
	sOne.SetRepoID(1)
	sOne.SetBuildID(1)
	sOne.SetNumber(1)
	sOne.SetName("foo")
	sOne.SetImage("baz")
	sOne.SetStatus("success")

	sTwo := testStep()
	sTwo.SetID(2)
	sTwo.SetRepoID(2)
	sTwo.SetBuildID(2)
	sTwo.SetNumber(1)
	sTwo.SetName("foo")
	sTwo.SetImage("baz")
	sTwo.SetStatus("success")

	sThree := testStep()
	sThree.SetID(3)
	sThree.SetRepoID(2)
	sThree.SetBuildID(2)
	sThree.SetNumber(2)
	sThree.SetName("bar")
	sThree.SetImage("bazian:latest")
	sThree.SetStatus("failure")

	want := make(map[string]float64)
	want["success"] = 2
	want["failure"] = 1

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()
	_ = db.CreateStep(sOne)
	_ = db.CreateStep(sTwo)
	_ = db.CreateStep(sThree)

	// run test
	got, err := db.GetStepStatusCount()

	if err != nil {
		t.Errorf("GetStepStatusCount returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetStepStatusCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateStep(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	want := testStep()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetName("foo")
	want.SetImage("baz")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateStep(want)

	if err != nil {
		t.Errorf("CreateStep returned err: %v", err)
	}

	got, _ := db.GetStep(int(want.GetNumber()), b)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateStep is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateStep_Invalid(t *testing.T) {
	// setup types
	s := testStep()
	s.SetID(1)
	s.SetBuildID(1)
	s.SetNumber(1)
	s.SetName("foo")
	s.SetImage("baz")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateStep(s)

	if err == nil {
		t.Errorf("CreateStep should have returned err")
	}
}

func TestDatabase_Client_UpdateStep(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	want := testStep()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetName("foo")
	want.SetImage("baz")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()
	_ = db.CreateStep(want)

	// run test
	err := db.UpdateStep(want)

	if err != nil {
		t.Errorf("UpdateStep returned err: %v", err)
	}

	got, _ := db.GetStep(int(want.GetNumber()), b)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateStep is %v, want %v", got, want)
	}
}

func TestDatabase_Client_UpdateStep_Invalid(t *testing.T) {
	// setup types
	s := testStep()
	s.SetID(1)
	s.SetBuildID(1)
	s.SetNumber(1)
	s.SetName("foo")
	s.SetImage("baz")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()
	_ = db.CreateStep(s)

	// run test
	err := db.UpdateStep(s)

	if err == nil {
		t.Errorf("UpdateStep should have returned err")
	}
}

func TestDatabase_Client_DeleteStep(t *testing.T) {
	// setup types
	want := testStep()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetName("foo")
	want.SetImage("baz")

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()
	_ = db.CreateStep(want)

	// run test
	err := db.DeleteStep(want.GetID())

	if err != nil {
		t.Errorf("DeleteStep returned err: %v", err)
	}
}

// testStep is a test helper function to create a
// library Step type with all fields set to their
// zero values.
func testStep() *library.Step {
	i64 := int64(0)
	i := 0
	str := ""
	return &library.Step{
		ID:           &i64,
		BuildID:      &i64,
		RepoID:       &i64,
		Number:       &i,
		Name:         &str,
		Image:        &str,
		Stage:        &str,
		Status:       &str,
		Error:        &str,
		ExitCode:     &i,
		Created:      &i64,
		Started:      &i64,
		Finished:     &i64,
		Host:         &str,
		Runtime:      &str,
		Distribution: &str,
	}
}
