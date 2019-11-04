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

	_, err = db.Database.DB().Exec(db.DDL.LogService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableLog, err)
	}
}

func TestDatabase_Client_GetBuildLogs(t *testing.T) {
	// setup types
	one := int64(1)
	data := []byte("foo")
	lOne := testLog()
	lOne.ID = &one
	lOne.StepID = &one
	lOne.BuildID = &one
	lOne.RepoID = &one
	lOne.Data = &data
	two := int64(2)
	lTwo := testLog()
	lTwo.ID = &two
	lTwo.StepID = &two
	lTwo.BuildID = &one
	lTwo.RepoID = &one
	lTwo.Data = &data
	want := []*library.Log{lOne, lTwo}

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from logs;")
		database.Database.Close()
	}()
	_ = database.CreateLog(lOne)
	_ = database.CreateLog(lTwo)

	// run test
	got, err := database.GetBuildLogs(1)

	if err != nil {
		t.Errorf("GetBuildLogs returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBuildLogs is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetStepLog(t *testing.T) {
	// setup types
	one := int64(1)
	data := []byte("foo")
	want := testLog()
	want.ID = &one
	want.BuildID = &one
	want.RepoID = &one
	want.ServiceID = &one
	want.StepID = &one
	want.Data = &data

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from logs;")
		database.Database.Close()
	}()
	_ = database.CreateLog(want)

	// run test
	got, err := database.GetStepLog(want.GetStepID())

	if err != nil {
		t.Errorf("GetLog returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetLog is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetServiceLog(t *testing.T) {
	// setup types
	one := int64(1)
	data := []byte("foo")
	want := testLog()
	want.ID = &one
	want.BuildID = &one
	want.RepoID = &one
	want.ServiceID = &one
	want.StepID = &one
	want.Data = &data

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from logs;")
		database.Database.Close()
	}()
	_ = database.CreateLog(want)

	// run test
	got, err := database.GetServiceLog(want.GetServiceID())

	if err != nil {
		t.Errorf("GetLog returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetLog is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateLog(t *testing.T) {
	// setup types
	one := int64(1)
	data := []byte("foo")
	want := testLog()
	want.ID = &one
	want.StepID = &one
	want.BuildID = &one
	want.RepoID = &one
	want.Data = &data

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from logs;")
		database.Database.Close()
	}()

	// run test
	err := database.CreateLog(want)

	if err != nil {
		t.Errorf("CreateLog returned err: %v", err)
	}

	got, _ := database.GetStepLog(want.GetStepID())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateLog is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateLog_Invalid(t *testing.T) {
	// setup types
	one := int64(1)
	data := []byte("foo")
	l := testLog()
	l.ID = &one
	l.BuildID = &one
	l.RepoID = &one
	l.Data = &data

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from logs;")
		database.Database.Close()
	}()

	// run test
	err := database.CreateLog(l)

	if err == nil {
		t.Errorf("CreateLog should have returned err")
	}
}

func TestDatabase_Client_UpdateLog(t *testing.T) {
	// setup types
	one := int64(1)
	data := []byte("foo")
	want := testLog()
	want.ID = &one
	want.StepID = &one
	want.BuildID = &one
	want.RepoID = &one
	want.Data = &data

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from logs;")
		database.Database.Close()
	}()
	_ = database.CreateLog(want)

	// run test
	err := database.UpdateLog(want)

	if err != nil {
		t.Errorf("UpdateLog returned err: %v", err)
	}

	got, _ := database.GetStepLog(want.GetStepID())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateLog is %v, want %v", got, want)
	}
}

func TestDatabase_Client_UpdateLog_Invalid(t *testing.T) {
	// setup types
	one := int64(1)
	data := []byte("foo")
	l := testLog()
	l.ID = &one
	l.BuildID = &one
	l.RepoID = &one
	l.Data = &data

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from logs;")
		database.Database.Close()
	}()
	_ = database.CreateLog(l)

	// run test
	err := database.UpdateLog(l)

	if err == nil {
		t.Errorf("UpdateLog should have returned err")
	}
}

func TestDatabase_Client_DeleteLog(t *testing.T) {
	// setup types
	one := int64(1)
	data := []byte("foo")
	want := testLog()
	want.ID = &one
	want.StepID = &one
	want.BuildID = &one
	want.RepoID = &one
	want.Data = &data

	// setup database
	database, _ := NewTest()
	defer func() {
		database.Database.Exec("delete from logs;")
		database.Database.Close()
	}()
	_ = database.CreateLog(want)

	// run test
	err := database.DeleteLog(want.GetID())

	if err != nil {
		t.Errorf("DeleteLog returned err: %v", err)
	}
}

// testLog is a test helper function to create a
// library Log type with all fields set to their
// zero values.
func testLog() *library.Log {
	i64 := int64(0)
	b := []byte{}
	return &library.Log{
		ID:        &i64,
		BuildID:   &i64,
		RepoID:    &i64,
		ServiceID: &i64,
		StepID:    &i64,
		Data:      &b,
	}
}
