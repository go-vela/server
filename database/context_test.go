// SPDX-License-Identifier: Apache-2.0

package database

import (
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/tracing"
)

func TestDatabase_FromContext(t *testing.T) {
	_postgres, _ := testPostgres(t)
	defer _postgres.Close()

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	ctx.Set(key, _postgres)

	typeCtx, _ := gin.CreateTestContext(nil)
	typeCtx.Set(key, nil)

	nilCtx, _ := gin.CreateTestContext(nil)
	nilCtx.Set(key, nil)

	// setup tests
	tests := []struct {
		name    string
		context *gin.Context
		want    Interface
	}{
		{
			name:    "success",
			context: ctx,
			want:    _postgres,
		},
		{
			name:    "failure with nil",
			context: nilCtx,
			want:    nil,
		},
		{
			name:    "failure with wrong type",
			context: typeCtx,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := FromContext(test.context)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("FromContext for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}

func TestDatabase_ToContext(t *testing.T) {
	context, _ := gin.CreateTestContext(nil)

	_postgres, _ := testPostgres(t)
	defer _postgres.Close()

	_sqlite := testSqlite(t)
	defer _sqlite.Close()

	// setup tests
	tests := []struct {
		name     string
		database *engine
		want     *engine
	}{
		{
			name:     "success with postgres",
			database: _postgres,
			want:     _postgres,
		},
		{
			name:     "success with sqlite3",
			database: _sqlite,
			want:     _sqlite,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ToContext(context, test.want)

			got := context.Value(key)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ToContext for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}

func TestDatabase_FromCLICommand(t *testing.T) {
	happyPath := new(cli.Command)

	happyPath.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "database.addr",
			Value: "file::memory:?cache=shared",
		},
		&cli.StringFlag{
			Name:  "database.driver",
			Value: "sqlite3",
		},
		&cli.IntFlag{
			Name:  "database.compression.level",
			Value: 3,
		},
		&cli.DurationFlag{
			Name:  "database.connection.life",
			Value: 10 * time.Second,
		},
		&cli.IntFlag{
			Name:  "database.connection.idle",
			Value: 5,
		},
		&cli.IntFlag{
			Name:  "database.connection.open",
			Value: 20,
		},
		&cli.StringFlag{
			Name:  "database.encryption.key",
			Value: "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
		},
		&cli.BoolFlag{
			Name:  "database.skip_creation",
			Value: true,
		},
	}

	// setup tests
	tests := []struct {
		name    string
		failure bool
		command *cli.Command
	}{
		{
			name:    "success",
			failure: false,
			command: happyPath,
		},
		{
			name:    "failure",
			failure: true,
			command: new(cli.Command),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := FromCLICommand(test.command, &tracing.Client{Config: tracing.Config{EnableTracing: false}})

			if test.failure {
				if err == nil {
					t.Errorf("FromCLICommand for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("FromCLICommand for %s returned err: %v", test.name, err)
			}
		})
	}
}
