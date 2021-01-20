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

	_, err = db.Database.DB().Exec(db.DDL.WorkerService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableWorker, err)
	}
}

func TestDatabase_Client_GetWorker(t *testing.T) {
	// setup types
	want := testWorker()
	want.SetID(1)
	want.SetHostname("worker_0")
	want.SetAddress("localhost")
	want.SetActive(true)
	want.SetBuildLimit(0)

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from workers;")
		db.Database.Close()
	}()

	_ = db.CreateWorker(want)

	// run test
	got, err := db.GetWorker(want.GetHostname())

	if err != nil {
		t.Errorf("GetWorker returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetWorker is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetWorkerList(t *testing.T) {
	// setup types
	wOne := testWorker()
	wOne.SetID(1)
	wOne.SetHostname("worker_1")
	wOne.SetAddress("localhost")
	wOne.SetActive(true)

	wTwo := testWorker()
	wTwo.SetID(2)
	wTwo.SetHostname("worker_2")
	wTwo.SetAddress("localhost")
	wTwo.SetActive(true)
	wTwo.SetBuildLimit(0)

	want := []*library.Worker{wOne, wTwo}

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from workers;")
		db.Database.Close()
	}()

	_ = db.CreateWorker(wOne)
	_ = db.CreateWorker(wTwo)

	// run test
	got, err := db.GetWorkerList()

	if err != nil {
		t.Errorf("GetWorkerList returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetWorkerList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetWorkerCount(t *testing.T) {
	// setup types
	wOne := testWorker()
	wOne.SetID(1)
	wOne.SetHostname("worker_1")
	wOne.SetAddress("localhost")
	wOne.SetActive(true)

	wTwo := testWorker()
	wTwo.SetID(2)
	wTwo.SetHostname("worker_2")
	wTwo.SetAddress("localhost")
	wTwo.SetActive(true)

	want := 2

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from workers;")
		db.Database.Close()
	}()

	_ = db.CreateWorker(wOne)
	_ = db.CreateWorker(wTwo)

	// run test
	got, err := db.GetWorkerCount()

	if err != nil {
		t.Errorf("GetWorkerCount returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("GetWorkerCount is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateWorker(t *testing.T) {
	// setup types
	want := testWorker()
	want.SetID(1)
	want.SetHostname("worker_0")
	want.SetAddress("localhost")
	want.SetActive(true)

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from workers;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateWorker(want)

	if err != nil {
		t.Errorf("CreateWorker returned err: %v", err)
	}

	got, _ := db.GetWorker(want.GetHostname())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateWorker is %v, want %v", got, want)
	}
}

func TestDatabase_Client_CreateWorker_Invalid(t *testing.T) {
	// setup types
	w := testWorker()
	w.SetID(1)
	w.SetAddress("localhost")
	w.SetActive(true)

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from workers;")
		db.Database.Close()
	}()

	// run test
	err := db.CreateWorker(w)

	if err == nil {
		t.Errorf("CreateWorker should have returned err")
	}
}

func TestDatabase_Client_UpdateWorker(t *testing.T) {
	// setup types
	want := testWorker()
	want.SetID(1)
	want.SetHostname("worker_0")
	want.SetAddress("localhost")
	want.SetActive(true)

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from workers;")
		db.Database.Close()
	}()

	_ = db.CreateWorker(want)

	// run test
	err := db.UpdateWorker(want)

	if err != nil {
		t.Errorf("UpdateWorker returned err: %v", err)
	}

	got, _ := db.GetWorker(want.GetHostname())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateWorker is %v, want %v", got, want)
	}
}

func TestDatabase_Client_UpdateWorker_Invalid(t *testing.T) {
	// setup types
	w := testWorker()
	w.SetID(1)
	w.SetHostname("worker_0")
	w.SetActive(true)

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from workers;")
		db.Database.Close()
	}()

	_ = db.CreateWorker(w)

	// run test
	err := db.UpdateWorker(w)

	if err == nil {
		t.Errorf("UpdateWorker should have returned err")
	}
}

func TestDatabase_Client_UpdateWorker_Boolean(t *testing.T) {
	// setup types
	want := testWorker()
	want.SetID(1)
	want.SetHostname("worker_0")
	want.SetAddress("localhost")
	want.SetActive(true)

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from workers;")
		db.Database.Close()
	}()

	_ = db.CreateWorker(want)

	// run test
	want.SetActive(false)

	err := db.UpdateWorker(want)

	if err != nil {
		t.Errorf("UpdateWorker returned err: %+v", err)
	}

	got, _ := db.GetWorker(want.GetHostname())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateWorker is %+v, want %+v", got, want)
	}
}

func TestDatabase_Client_DeleteWorker(t *testing.T) {
	// setup types
	want := testWorker()
	want.SetID(1)
	want.SetHostname("worker_0")
	want.SetAddress("localhost")
	want.SetActive(true)

	// setup database
	db, _ := NewTest()

	defer func() {
		db.Database.Exec("delete from workers;")
		db.Database.Close()
	}()

	_ = db.CreateWorker(want)

	// run test
	err := db.DeleteWorker(want.GetID())

	if err != nil {
		t.Errorf("DeleteWorker returned err: %v", err)
	}
}

// testWorker is a test helper function to create a
// library Worker type with all fields set to their
// zero values.
func testWorker() *library.Worker {
	i64 := int64(0)
	str := ""
	b := false
	arr := []string{}

	return &library.Worker{
		ID:            &i64,
		Hostname:      &str,
		Address:       &str,
		Routes:        &arr,
		Active:        &b,
		LastCheckedIn: &i64,
		BuildLimit:    &i64,
	}
}
