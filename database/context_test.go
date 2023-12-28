// SPDX-License-Identifier: Apache-2.0

package database

import (
	"flag"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
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

func TestDatabase_FromCLIContext(t *testing.T) {
	flags := flag.NewFlagSet("test", 0)
	flags.String("database.driver", "sqlite3", "doc")
	flags.String("database.addr", "file::memory:?cache=shared", "doc")
	flags.Int("database.compression.level", 3, "doc")
	flags.Duration("database.connection.life", 10*time.Second, "doc")
	flags.Int("database.connection.idle", 5, "doc")
	flags.Int("database.connection.open", 20, "doc")
	flags.String("database.encryption.key", "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW", "doc")
	flags.Bool("database.skip_creation", true, "doc")

	// setup tests
	tests := []struct {
		name    string
		failure bool
		context *cli.Context
	}{
		{
			name:    "success",
			failure: false,
			context: cli.NewContext(&cli.App{Name: "vela"}, flags, nil),
		},
		{
			name:    "failure",
			failure: true,
			context: cli.NewContext(&cli.App{Name: "vela"}, flag.NewFlagSet("test", 0), nil),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := FromCLIContext(test.context)

			if test.failure {
				if err == nil {
					t.Errorf("FromCLIContext for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("FromCLIContext for %s returned err: %v", test.name, err)
			}
		})
	}
}
