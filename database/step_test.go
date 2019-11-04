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
	one := int64(1)
	oneNum := 1
	oneName := "foo"
	b := testBuild()
	b.ID = &one
	b.RepoID = &one
	b.Number = &oneNum
	want := testStep()
	want.ID = &one
	want.RepoID = &one
	want.BuildID = &one
	want.Number = &oneNum
	want.Name = &oneName

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
	one := int64(1)
	oneNum := 1
	oneName := "foo"
	sOne := testStep()
	sOne.ID = &one
	sOne.RepoID = &one
	sOne.BuildID = &one
	sOne.Number = &oneNum
	sOne.Name = &oneName
	two := int64(2)
	twoNum := 2
	twoName := "bar"
	sTwo := testStep()
	sTwo.ID = &two
	sTwo.RepoID = &one
	sTwo.BuildID = &one
	sTwo.Number = &twoNum
	sTwo.Name = &twoName
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
	one := int64(1)
	oneNum := 1
	oneName := "foo"
	b := testBuild()
	b.ID = &one
	b.RepoID = &one
	b.Number = &oneNum
	sOne := testStep()
	sOne.ID = &one
	sOne.RepoID = &one
	sOne.BuildID = &one
	sOne.Number = &oneNum
	sOne.Name = &oneName
	two := int64(2)
	twoNum := 2
	twoName := "bar"
	sTwo := testStep()
	sTwo.ID = &two
	sTwo.RepoID = &one
	sTwo.BuildID = &one
	sTwo.Number = &twoNum
	sTwo.Name = &twoName
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
	one := int64(1)
	oneNum := 1
	oneName := "foo"
	b := testBuild()
	b.ID = &one
	b.RepoID = &one
	b.Number = &oneNum
	sOne := testStep()
	sOne.ID = &one
	sOne.RepoID = &one
	sOne.BuildID = &one
	sOne.Number = &oneNum
	sOne.Name = &oneName
	two := int64(2)
	twoNum := 2
	twoName := "bar"
	sTwo := testStep()
	sTwo.ID = &two
	sTwo.RepoID = &one
	sTwo.BuildID = &one
	sTwo.Number = &twoNum
	sTwo.Name = &twoName
	three := int64(3)
	threeNum := 3
	threeName := "baz"
	sThree := testStep()
	sThree.ID = &three
	sThree.RepoID = &one
	sThree.BuildID = &two
	sThree.Number = &threeNum
	sThree.Name = &threeName

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

func TestDatabase_Client_CreateStep(t *testing.T) {
	// setup types
	one := int64(1)
	oneNum := 1
	oneName := "foo"
	b := testBuild()
	b.ID = &one
	b.RepoID = &one
	b.Number = &oneNum
	want := testStep()
	want.ID = &one
	want.RepoID = &one
	want.BuildID = &one
	want.Number = &oneNum
	want.Name = &oneName

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
	one := int64(1)
	oneNum := 1
	oneName := "foo"
	b := testBuild()
	b.ID = &one
	b.RepoID = &one
	b.Number = &oneNum
	s := testStep()
	s.ID = &one
	s.BuildID = &one
	s.Number = &oneNum
	s.Name = &oneName

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
	one := int64(1)
	oneNum := 1
	oneName := "foo"
	b := testBuild()
	b.ID = &one
	b.RepoID = &one
	b.Number = &oneNum
	want := testStep()
	want.ID = &one
	want.RepoID = &one
	want.BuildID = &one
	want.Number = &oneNum
	want.Name = &oneName

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
	one := int64(1)
	oneNum := 1
	oneName := "foo"
	b := testBuild()
	b.ID = &one
	b.RepoID = &one
	b.Number = &oneNum
	s := testStep()
	s.ID = &one
	s.BuildID = &one
	s.Number = &oneNum
	s.Name = &oneName

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
	one := int64(1)
	oneNum := 1
	oneName := "foo"
	want := testStep()
	want.ID = &one
	want.RepoID = &one
	want.BuildID = &one
	want.Number = &oneNum
	want.Name = &oneName

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
