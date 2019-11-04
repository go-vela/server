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
	"github.com/google/go-cmp/cmp"
	"github.com/kr/pretty"
)

func init() {
	db, err := NewTest()
	if err != nil {
		log.Fatalf("Error creating test database: %v", err)
	}

	_, _ = db.Database.DB().Exec(db.DDL.ServiceService.Create)
	if err != nil {
		log.Fatalf("Error creating %s table: %v", constants.TableService, err)
	}
}

func TestDatabase_Client_GetService(t *testing.T) {
	// setup types
	zero := int64(0)
	one := 1
	one64 := int64(1)
	foo := "foo"
	b := &library.Build{
		ID:     &one64,
		RepoID: &one64,
		Number: &one,
	}
	want := &library.Service{
		ID:       &one64,
		RepoID:   &one64,
		BuildID:  &one64,
		Created:  &zero,
		Started:  &zero,
		Finished: &zero,
		Number:   &one,
		Name:     &foo,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()
	err := db.CreateService(want)

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
	one := 1
	one64 := int64(1)
	two64 := int64(2)
	foo := "foo"
	sOne := &library.Service{
		ID:       &one64,
		RepoID:   &one64,
		BuildID:  &one64,
		Created:  &one64,
		Started:  &one64,
		Finished: &one64,
		Number:   &one,
		Name:     &foo,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}
	sTwo := &library.Service{
		ID:       &two64,
		RepoID:   &one64,
		BuildID:  &two64,
		Created:  &one64,
		Started:  &one64,
		Finished: &one64,
		Number:   &one,
		Name:     &foo,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}
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

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetServiceList() mismatch (-want +got):\n%s", diff)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetServiceList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetBuildServiceList(t *testing.T) {
	// setup types
	one := 1
	two := 2
	one64 := int64(1)
	two64 := int64(2)
	foo := "foo"
	bar := "bar"

	b := &library.Build{
		ID:     &one64,
		RepoID: &one64,
		Number: &one,
	}
	sOne := &library.Service{
		ID:       &one64,
		RepoID:   &one64,
		BuildID:  &one64,
		Created:  &one64,
		Started:  &one64,
		Finished: &one64,
		Number:   &one,
		Name:     &foo,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}
	sTwo := &library.Service{
		ID:       &two64,
		RepoID:   &one64,
		BuildID:  &one64,
		Created:  &one64,
		Started:  &one64,
		Finished: &one64,
		Number:   &two,
		Name:     &bar,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}
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
		pretty.Ldiff(t, got, want)
		t.Errorf("GetBuildServiceList is %v, want %v", got, want)
	}
}

func TestDatabase_Client_GetBuildServiceCount(t *testing.T) {
	// setup types
	one := 1
	two := 2
	one64 := int64(1)
	two64 := int64(2)
	three64 := int64(3)
	foo := "foo"
	bar := "bar"

	b := &library.Build{
		ID:     &two64,
		RepoID: &one64,
		Number: &two,
	}
	sOne := &library.Service{
		ID:       &one64,
		RepoID:   &one64,
		BuildID:  &one64,
		Created:  &one64,
		Started:  &one64,
		Finished: &one64,
		Number:   &one,
		Name:     &foo,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}
	sTwo := &library.Service{
		ID:       &two64,
		RepoID:   &two64,
		BuildID:  &two64,
		Created:  &one64,
		Started:  &one64,
		Finished: &one64,
		Number:   &one,
		Name:     &bar,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}
	sThree := &library.Service{
		ID:       &three64,
		RepoID:   &two64,
		BuildID:  &two64,
		Created:  &one64,
		Started:  &one64,
		Finished: &one64,
		Number:   &two,
		Name:     &bar,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}
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

func TestDatabase_Client_CreateService(t *testing.T) {
	// setup types
	zero := int64(0)
	one := 1
	one64 := int64(1)
	foo := "foo"
	b := &library.Build{
		ID:     &one64,
		RepoID: &one64,
		Number: &one,
	}
	want := &library.Service{
		ID:       &one64,
		RepoID:   &one64,
		BuildID:  &one64,
		Created:  &zero,
		Started:  &zero,
		Finished: &zero,
		Number:   &one,
		Name:     &foo,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}

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
	zero := int64(0)
	one := 1
	one64 := int64(1)
	foo := "foo"
	s := &library.Service{
		ID:       &one64,
		RepoID:   &one64,
		BuildID:  &one64,
		Created:  &zero,
		Started:  &zero,
		Finished: &zero,
		Number:   &one,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}

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
	zero := int64(0)
	one := 1
	one64 := int64(1)
	foo := "foo"
	b := &library.Build{
		ID:     &one64,
		RepoID: &one64,
		Number: &one,
	}
	want := &library.Service{
		ID:       &one64,
		RepoID:   &one64,
		BuildID:  &one64,
		Created:  &zero,
		Started:  &zero,
		Finished: &zero,
		Number:   &one,
		Name:     &foo,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}

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
	zero := int64(0)
	one := 1
	one64 := int64(1)
	foo := "foo"
	s := &library.Service{
		ID:       &one64,
		RepoID:   &one64,
		BuildID:  &one64,
		Created:  &zero,
		Started:  &zero,
		Finished: &zero,
		Number:   &one,
		Name:     &foo,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()
	_ = db.CreateService(s)

	// run test
	s.RepoID = &zero
	err := db.UpdateService(s)

	if err == nil {
		t.Errorf("UpdateService should have returned err")
	}
}

func TestDatabase_Client_DeleteService(t *testing.T) {
	// setup types
	zero := int64(0)
	one := 1
	one64 := int64(1)
	foo := "foo"
	want := &library.Service{
		ID:       &one64,
		RepoID:   &one64,
		BuildID:  &one64,
		Created:  &zero,
		Started:  &zero,
		Finished: &zero,
		Number:   &one,
		Name:     &foo,
		Status:   &foo,
		Error:    &foo,
		ExitCode: &one,
	}

	// setup database
	db, _ := NewTest()
	defer func() {
		db.Database.Exec("delete from services;")
		db.Database.Close()
	}()
	_ = db.CreateService(want)

	// run test
	err := db.DeleteService(want.GetBuildID())

	if err != nil {
		t.Errorf("DeleteService returned err: %v", err)
	}
}
