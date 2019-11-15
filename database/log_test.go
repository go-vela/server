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
	lOne := testLog()
	lOne.SetID(1)
	lOne.SetStepID(1)
	lOne.SetBuildID(1)
	lOne.SetRepoID(1)
	lOne.SetData([]byte{})

	lTwo := testLog()
	lTwo.SetID(2)
	lTwo.SetStepID(2)
	lTwo.SetBuildID(1)
	lTwo.SetRepoID(1)
	lTwo.SetData([]byte{})

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
	want := testLog()
	want.SetID(1)
	want.SetBuildID(1)
	want.SetRepoID(1)
	want.SetStepID(1)
	want.SetData([]byte{})

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
	want := testLog()
	want.SetID(1)
	want.SetBuildID(1)
	want.SetRepoID(1)
	want.SetServiceID(1)
	want.SetData([]byte{})

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
	want := testLog()
	want.SetID(1)
	want.SetBuildID(1)
	want.SetRepoID(1)
	want.SetStepID(1)
	want.SetData([]byte{})

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
	l := testLog()
	l.SetID(1)
	l.SetBuildID(1)
	l.SetRepoID(1)
	l.SetData([]byte{})

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
	want := testLog()
	want.SetID(1)
	want.SetStepID(1)
	want.SetBuildID(1)
	want.SetRepoID(1)
	want.SetData([]byte{})

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
	l := testLog()
	l.SetID(1)
	l.SetBuildID(1)
	l.SetRepoID(1)
	l.SetData([]byte{})

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
	want := testLog()
	want.SetID(1)
	want.SetStepID(1)
	want.SetBuildID(1)
	want.SetRepoID(1)
	want.SetData([]byte{})

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
