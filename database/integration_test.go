// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-vela/server/database/user"

	"github.com/go-vela/server/database/worker"

	"github.com/go-vela/types/library"
)

func TestDatabase_Integration(t *testing.T) {
	// check if we should skip the integration test
	if os.Getenv("INTEGRATION") == "" {
		t.Skipf("skipping %s integration test due to environment variable constraint", t.Name())
	}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		config  *config
	}{
		{
			name:    "success with postgres",
			failure: false,
			config: &config{
				Driver:           "postgres",
				Address:          "postgres://vela:notARealPassword12345@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, err := New(
				WithAddress(test.config.Address),
				WithCompressionLevel(test.config.CompressionLevel),
				WithConnectionLife(test.config.ConnectionLife),
				WithConnectionIdle(test.config.ConnectionIdle),
				WithConnectionOpen(test.config.ConnectionOpen),
				WithDriver(test.config.Driver),
				WithEncryptionKey(test.config.EncryptionKey),
				WithSkipCreation(test.config.SkipCreation),
			)
			if err != nil {
				t.Errorf("unable to create new database engine for %s: %v", test.name, err)
			}

			driver := db.Driver()
			if !strings.EqualFold(driver, test.config.Driver) {
				t.Errorf("Driver() is %v, want %v", driver, test.config.Driver)
			}

			err = db.Ping()
			if err != nil {
				t.Errorf("unable to ping database engine for %s: %v", test.name, err)
			}

			t.Run("test_users", func(t *testing.T) {
				testUsers(t, db)
			})

			t.Run("test_workers", func(t *testing.T) {
				testWorkers(t, db)
			})

			err = db.Close()
			if err != nil {
				t.Errorf("unable to close database engine for %s: %v", test.name, err)
			}
		})
	}
}

func testUsers(t *testing.T, db Interface) {
	// used to track the number of methods we call for users
	counter := 0

	one := new(library.User)
	one.SetID(1)
	one.SetName("octocat")
	one.SetToken("superSecretToken")
	one.SetRefreshToken("superSecretRefreshToken")
	one.SetHash("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy")
	one.SetFavorites([]string{"github/octocat"})
	one.SetActive(true)
	one.SetAdmin(false)

	two := new(library.User)
	two.SetID(2)
	two.SetName("octocat")
	two.SetToken("superSecretToken")
	two.SetRefreshToken("superSecretRefreshToken")
	two.SetHash("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy")
	two.SetFavorites([]string{"github/octocat"})
	two.SetActive(true)
	two.SetAdmin(false)

	users := []*library.User{one, two}

	// create the users
	for _, user := range users {
		err := db.CreateUser(user)
		if err != nil {
			t.Errorf("unable to create user %s: %v", user.GetName(), err)
		}
	}
	counter++

	// count the users
	count, err := db.CountUsers()
	if err != nil {
		t.Errorf("unable to count users: %v", err)
	}
	if int(count) != len(users) {
		t.Errorf("CountUsers() is %v, want 2", count)
	}
	counter++

	// list the users
	list, err := db.ListUsers()
	if err != nil {
		t.Errorf("unable to list users: %v", err)
	}
	if !reflect.DeepEqual(list, users) {
		t.Errorf("ListUsers() is %v, want %v", list, users)
	}
	counter++

	// list the users
	lite, count, err := db.ListLiteUsers(1, 10)
	if err != nil {
		t.Errorf("unable to list lite users: %v", err)
	}
	if !reflect.DeepEqual(lite, users) {
		t.Errorf("ListLiteUsers() is %v, want %v", list, users)
	}
	if int(count) != len(users) {
		t.Errorf("ListLiteUsers() is %v, want %v", count, len(users))
	}
	counter++

	// lookup the users by name
	for _, user := range users {
		got, err := db.GetUserForName(user.GetName())
		if err != nil {
			t.Errorf("unable to get user %s by hostname: %v", user.GetName(), err)
		}
		if !reflect.DeepEqual(got, user) {
			t.Errorf("GetUserForName() is %v, want %v", got, user)
		}
	}
	counter++

	// update the users
	for _, user := range users {
		user.SetActive(false)
		err = db.UpdateUser(user)
		if err != nil {
			t.Errorf("unable to update user %s: %v", user.GetName(), err)
		}

		// lookup the user by ID
		got, err := db.GetUser(user.GetID())
		if err != nil {
			t.Errorf("unable to get user %s by ID: %v", user.GetName(), err)
		}
		if !reflect.DeepEqual(got, user) {
			t.Errorf("GetUser() is %v, want %v", got, user)
		}
	}
	counter++
	counter++

	// delete the users
	for _, user := range users {
		err = db.DeleteUser(user)
		if err != nil {
			t.Errorf("unable to delete user %s: %v", user.GetName(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	//
	// we subtract 2 for creating the table and indexes for users
	// since those are already called when the database engine starts
	methods := reflect.TypeOf(new(user.UserInterface)).Elem().NumMethod() - 2
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testWorkers(t *testing.T, db Interface) {
	// used to track the number of methods we call for workers
	counter := 0

	one := new(library.Worker)
	one.SetID(1)
	one.SetHostname("worker-1.example.com")
	one.SetAddress("https://worker-1.example.com")
	one.SetRoutes([]string{"vela"})
	one.SetActive(true)
	one.SetStatus("available")
	one.SetLastStatusUpdateAt(time.Now().UTC().Unix())
	one.SetRunningBuildIDs([]string{"12345"})
	one.SetLastBuildStartedAt(time.Now().UTC().Unix())
	one.SetLastBuildFinishedAt(time.Now().UTC().Unix())
	one.SetLastCheckedIn(time.Now().UTC().Unix())
	one.SetBuildLimit(1)

	two := new(library.Worker)
	two.SetID(2)
	two.SetHostname("worker-2.example.com")
	two.SetAddress("https://worker-2.example.com")
	two.SetRoutes([]string{"vela"})
	two.SetActive(true)
	two.SetStatus("available")
	two.SetLastStatusUpdateAt(time.Now().UTC().Unix())
	two.SetRunningBuildIDs([]string{"12345"})
	two.SetLastBuildStartedAt(time.Now().UTC().Unix())
	two.SetLastBuildFinishedAt(time.Now().UTC().Unix())
	two.SetLastCheckedIn(time.Now().UTC().Unix())
	two.SetBuildLimit(1)

	workers := []*library.Worker{one, two}

	// create the workers
	for _, worker := range workers {
		err := db.CreateWorker(worker)
		if err != nil {
			t.Errorf("unable to create worker %s: %v", worker.GetHostname(), err)
		}
	}
	counter++

	// count the workers
	count, err := db.CountWorkers()
	if err != nil {
		t.Errorf("unable to count workers: %v", err)
	}
	if int(count) != len(workers) {
		t.Errorf("CountWorkers() is %v, want %v", count, len(workers))
	}
	counter++

	// list the workers
	list, err := db.ListWorkers()
	if err != nil {
		t.Errorf("unable to list workers: %v", err)
	}
	if !reflect.DeepEqual(list, workers) {
		t.Errorf("ListWorkers() is %v, want %v", list, workers)
	}
	counter++

	// lookup the workers by hostname
	for _, worker := range workers {
		got, err := db.GetWorkerForHostname(worker.GetHostname())
		if err != nil {
			t.Errorf("unable to get worker %s by hostname: %v", worker.GetHostname(), err)
		}
		if !reflect.DeepEqual(got, worker) {
			t.Errorf("GetWorkerForHostname() is %v, want %v", got, worker)
		}
	}
	counter++

	// update the workers
	for _, worker := range workers {
		worker.SetActive(false)
		err = db.UpdateWorker(worker)
		if err != nil {
			t.Errorf("unable to update worker %s: %v", worker.GetHostname(), err)
		}

		// lookup the worker by ID
		got, err := db.GetWorker(worker.GetID())
		if err != nil {
			t.Errorf("unable to get worker %s by ID: %v", worker.GetHostname(), err)
		}
		if !reflect.DeepEqual(got, worker) {
			t.Errorf("GetWorker() is %v, want %v", got, worker)
		}
	}
	counter++
	counter++

	// delete the workers
	for _, worker := range workers {
		err = db.DeleteWorker(worker)
		if err != nil {
			t.Errorf("unable to delete worker %s: %v", worker.GetHostname(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	//
	// we subtract 2 for creating the table and indexes for workers
	// since those are already called when the database engine starts
	methods := reflect.TypeOf(new(worker.WorkerInterface)).Elem().NumMethod() - 2
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}
