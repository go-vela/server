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

	_, err = db.Database.DB().Exec(db.DDL.ServiceService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableService, err)
	}
}

func TestDatabase_Client_GetService(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	want := testService()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetName("foo")
	want.SetImage("bar")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()

	_ = db.CreateService(want)

	// run test
	got, err := db.GetService(want.GetNumber(), b)

	if err != nil {
		t.Errorf("GetService returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetService is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetServiceList(t *testing.T) {
	// setup types
	sOne := testService()
	sOne.SetID(1)
	sOne.SetRepoID(1)
	sOne.SetBuildID(1)
	sOne.SetNumber(1)
	sOne.SetName("foo")
	sOne.SetImage("baz")

	sTwo := testService()
	sTwo.SetID(2)
	sTwo.SetRepoID(1)
	sTwo.SetBuildID(1)
	sTwo.SetNumber(2)
	sTwo.SetName("bar")
	sTwo.SetImage("baz")

	want := []*library.Service{sOne, sTwo}

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()

	_ = db.CreateService(sOne)
	_ = db.CreateService(sTwo)

	// run test
	got, err := db.GetServiceList()

	if err != nil {
		t.Errorf("GetServiceList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetServiceList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetBuildServiceList(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	sOne := testService()
	sOne.SetID(1)
	sOne.SetRepoID(1)
	sOne.SetBuildID(1)
	sOne.SetNumber(1)
	sOne.SetName("foo")
	sOne.SetImage("baz")

	sTwo := testService()
	sTwo.SetID(2)
	sTwo.SetRepoID(1)
	sTwo.SetBuildID(1)
	sTwo.SetNumber(2)
	sTwo.SetName("bar")
	sTwo.SetImage("baz")

	want := []*library.Service{sTwo, sOne}

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()

	_ = db.CreateService(sOne)
	_ = db.CreateService(sTwo)

	// run test
	got, err := db.GetBuildServiceList(b, 1, 10)

	if err != nil {
		t.Errorf("GetBuildServiceList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBuildServiceList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetBuildServiceCount(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(2)
	b.SetRepoID(1)
	b.SetNumber(1)

	sOne := testService()
	sOne.SetID(1)
	sOne.SetRepoID(1)
	sOne.SetBuildID(1)
	sOne.SetNumber(1)
	sOne.SetName("foo")
	sOne.SetImage("baz")

	sTwo := testService()
	sTwo.SetID(2)
	sTwo.SetRepoID(2)
	sTwo.SetBuildID(2)
	sTwo.SetNumber(1)
	sTwo.SetName("foo")
	sTwo.SetImage("baz")

	sThree := testService()
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
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()

	_ = db.CreateService(sOne)
	_ = db.CreateService(sTwo)
	_ = db.CreateService(sThree)

	// run test
	got, err := db.GetBuildServiceCount(b)

	if err != nil {
		t.Errorf("GetBuildServiceCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetBuildServiceCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetServiceImageCount(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(2)
	b.SetRepoID(1)
	b.SetNumber(1)

	sOne := testService()
	sOne.SetID(1)
	sOne.SetRepoID(1)
	sOne.SetBuildID(1)
	sOne.SetNumber(1)
	sOne.SetName("foo")
	sOne.SetImage("baz")

	sTwo := testService()
	sTwo.SetID(2)
	sTwo.SetRepoID(2)
	sTwo.SetBuildID(2)
	sTwo.SetNumber(1)
	sTwo.SetName("foo")
	sTwo.SetImage("baz")

	sThree := testService()
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
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()

	_ = db.CreateService(sOne)
	_ = db.CreateService(sTwo)
	_ = db.CreateService(sThree)

	// run test
	got, err := db.GetServiceImageCount()

	if err != nil {
		t.Errorf("GetServiceImageCount returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetServiceImageCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateService(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	want := testService()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetName("foo")
	want.SetImage("baz")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateService(want)

	if err != nil {
		t.Errorf("CreateService returned err: %v", err)
	}

	got, _ := db.GetService(want.GetNumber(), b)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateService is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateService_Invalid(t *testing.T) {
	// setup types
	s := testService()
	s.SetID(1)
	s.SetRepoID(1)
	s.SetBuildID(1)
	s.SetName("foo")
	s.SetImage("baz")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateService(s)

	if err == nil {
		t.Errorf("CreateService should have returned err")
	}
}

func TestDatabase_Client_UpdateService(t *testing.T) {
	// setup types
	b := testBuild()
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	want := testService()
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetName("foo")
	want.SetImage("baz")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()

	_ = db.CreateService(want)

	// run test
	err := db.UpdateService(want)

	if err != nil {
		t.Errorf("UpdateService returned err: %v", err)
	}

	got, _ := db.GetService(want.GetNumber(), b)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateService is %v, want %v", got, want)
	}
}

func TestDatabase_Client_UpdateService_Invalid(t *testing.T) {
	// setup types
	s := testService()
	s.SetID(1)
	s.SetRepoID(1)
	s.SetBuildID(1)
	s.SetName("foo")
	s.SetImage("foo")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()

	_ = db.CreateService(s)

	// run test
	err := db.UpdateService(s)

	if err == nil {
		t.Errorf("UpdateService should have returned err")
	}
}

func TestDatabase_Client_DeleteService(t *testing.T) {
	// setup types
	s := testService()
	s.SetID(1)
	s.SetRepoID(1)
	s.SetBuildID(1)
	s.SetNumber(1)
	s.SetName("foo")
	s.SetImage("image")

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()

	_ = db.CreateService(s)

	// run test
	err := db.DeleteService(s.GetBuildID())

	if err != nil {
		t.Errorf("DeleteService returned err: %v", err)
	}
}

// testService is a test helper function to create a
// library Service type with all fields set to their
// zero values.
func testService() *library.Service {
	i64 := int64(0)
	i := 0
	str := ""

	return &library.Service{
		ID:       &i64,
		BuildID:  &i64,
		RepoID:   &i64,
		Number:   &i,
		Name:     &str,
		Image:    &str,
		Status:   &str,
		Error:    &str,
		ExitCode: &i,
		Created:  &i64,
		Started:  &i64,
		Finished: &i64,
	}
}
