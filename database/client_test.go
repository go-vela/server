// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"os"
	"testing"
	"time"

	"github.com/go-vela/server/database/ddl"

	"github.com/go-vela/types/constants"

	"github.com/jinzhu/gorm"
)

func TestDatabase_New(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	s := &Setup{
		Driver:           name,
		Address:          config,
		CompressionLevel: 3,
		ConnectionIdle:   2,
		ConnectionLife:   30 * time.Minute,
		ConnectionOpen:   0,
		EncryptionKey:    "C639A572E14D5075C526FDDD43E4ECF6",
	}

	// run test
	database, err := New(s)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	// nolint: staticcheck // ignore false positive
	defer database.Database.Close()

	if database == nil {
		t.Error("New returned nil database")
	}
}

func TestDatabase_New_Empty(t *testing.T) {
	// setup types
	s := &Setup{
		Driver:           "",
		Address:          "",
		CompressionLevel: 3,
		ConnectionIdle:   2,
		ConnectionLife:   30 * time.Minute,
		ConnectionOpen:   0,
		EncryptionKey:    "C639A572E14D5075C526FDDD43E4ECF6",
	}

	// run test
	database, err := New(s)

	if err == nil {
		t.Errorf("New should have returned err")
	}

	if database != nil {
		t.Errorf("New is %v want nil", database)
	}
}

func TestDatabase_NewTest(t *testing.T) {
	// run test
	database, err := NewTest()
	if err != nil {
		t.Errorf("newTest returned err: %v", err)
	}

	// nolint: staticcheck // ignore false positive
	defer database.Database.Close()

	if database == nil {
		t.Error("newTest returned nil database")
	}
}

func TestDatabase_setupDatabase(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)

	// run test
	err := setupDatabase(database.DB(), ddlMap)
	if err != nil {
		t.Errorf("setupDatabase returned err: %v", err)
	}
}

func TestDatabase_setupDatabase_BadDatabase(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	ddlMap, _ := ddl.NewMap(name)

	// run test
	database.Close()

	err := setupDatabase(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("setupDatabase should have returned err")
	}
}

func TestDatabase_setupDatabase_BadTable(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)

	// run test
	ddlMap.BuildService.Create = "#"

	err := setupDatabase(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("setupDatabase should have returned err")
	}
}

func TestDatabase_setupDatabase_BadIndex(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)

	// run test
	ddlMap.BuildService.Indexes = []string{"#"}

	err := setupDatabase(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("setupDatabase should have returned err")
	}
}

func TestDatabase_pingDatabase(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	// run test
	err := pingDatabase(database.DB())
	if err != nil {
		t.Errorf("pingDatabase returned err: %v", err)
	}
}

func TestDatabase_pingDatabase_BadDatabase(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)

	// run test
	database.Close()

	err := pingDatabase(database.DB())
	if err == nil {
		t.Errorf("pingDatabase should have returned err")
	}
}

func TestDatabase_createTables(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)

	// run test
	err := createTables(database.DB(), ddlMap)
	if err != nil {
		t.Errorf("createTables returned err: %v", err)
	}
}

func TestDatabase_createTables_BadBuildTable(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)

	ddlMap, _ := ddl.NewMap(name)

	// run test
	ddlMap.BuildService.Create = "#"

	err := createTables(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createTables should have returned err")
	}
}

func TestDatabase_createTables_BadLogTable(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)

	ddlMap, _ := ddl.NewMap(name)

	// run test
	ddlMap.LogService.Create = "#"

	err := createTables(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createTables should have returned err")
	}
}

func TestDatabase_createTables_BadRepoTable(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)

	ddlMap, _ := ddl.NewMap(name)

	// run test
	ddlMap.RepoService.Create = "#"

	err := createTables(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createTables should have returned err")
	}
}

func TestDatabase_createTables_BadSecretTable(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)

	ddlMap, _ := ddl.NewMap(name)

	// run test
	ddlMap.SecretService.Create = "#"

	err := createTables(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createTables should have returned err")
	}
}

func TestDatabase_createTables_BadStepTable(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)

	ddlMap, _ := ddl.NewMap(name)

	// run test
	ddlMap.StepService.Create = "#"

	err := createTables(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createTables should have returned err")
	}
}

func TestDatabase_createTables_BadUserTable(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)

	ddlMap, _ := ddl.NewMap(name)

	// run test
	ddlMap.UserService.Create = "#"

	err := createTables(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createTables should have returned err")
	}
}

func TestDatabase_createTables_BadWorkerTable(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)

	ddlMap, _ := ddl.NewMap(name)

	// run test
	ddlMap.WorkerService.Create = "#"

	err := createTables(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createTables should have returned err")
	}
}

func TestDatabase_createIndexes(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)
	_ = createTables(database.DB(), ddlMap)

	// run test
	err := createIndexes(database.DB(), ddlMap)
	if err != nil {
		t.Errorf("createIndexes returned err: %v", err)
	}
}

func TestDatabase_createIndexes_BadBuildIndex(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)
	_ = createTables(database.DB(), ddlMap)

	// run test
	ddlMap.BuildService.Indexes = []string{"#"}

	err := createIndexes(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createIndexes should have returned err")
	}
}

func TestDatabase_createIndexes_BadLogIndex(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)
	_ = createTables(database.DB(), ddlMap)

	// run test
	ddlMap.LogService.Indexes = []string{"#"}

	err := createIndexes(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createIndexes should have returned err")
	}
}

func TestDatabase_createIndexes_BadRepoIndex(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)
	_ = createTables(database.DB(), ddlMap)

	// run test
	ddlMap.RepoService.Indexes = []string{"#"}

	err := createIndexes(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createIndexes should have returned err")
	}
}

func TestDatabase_createIndexes_BadSecretIndex(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)
	_ = createTables(database.DB(), ddlMap)

	// run test
	ddlMap.SecretService.Indexes = []string{"#"}

	err := createIndexes(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createIndexes should have returned err")
	}
}

func TestDatabase_createIndexes_BadStepIndex(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)
	_ = createTables(database.DB(), ddlMap)

	// run test
	ddlMap.StepService.Indexes = []string{"#"}

	err := createIndexes(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createIndexes should have returned err")
	}
}

func TestDatabase_createIndexes_BadUserIndex(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)
	_ = createTables(database.DB(), ddlMap)

	// run test
	ddlMap.UserService.Indexes = []string{"#"}

	err := createIndexes(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createIndexes should have returned err")
	}
}

func TestDatabase_createIndexes_BadWorkerIndex(t *testing.T) {
	// setup types
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// setup database
	database, _ := gorm.Open(name, config)
	defer database.Close()

	ddlMap, _ := ddl.NewMap(name)
	_ = createTables(database.DB(), ddlMap)

	// run test
	ddlMap.WorkerService.Indexes = []string{"#"}

	err := createIndexes(database.DB(), ddlMap)
	if err == nil {
		t.Errorf("createIndexes should have returned err")
	}
}
