// SPDX-License-Identifier: Apache-2.0

package constants

// Server database drivers.
const (
	// DriverPostgres defines the driver type when integrating with a PostgreSQL database.
	DriverPostgres = "postgres"

	// DriverSqlite defines the driver type when integrating with a SQLite database.
	DriverSqlite = "sqlite3"
)

// Worker executor drivers.
const (
	// DriverDarwin defines the driver type when integrating with a darwin distribution.
	DriverDarwin = "darwin"

	// DriverLinux defines the driver type when integrating with a linux distribution.
	DriverLinux = "linux"

	// DriverLocal defines the driver type when integrating with a local system.
	DriverLocal = "local"

	// DriverWindows defines the driver type when integrating with a windows distribution.
	DriverWindows = "windows"
)

// Server and worker queue drivers.
const (

	// DriverKafka defines the driver type when integrating with a Kafka queue.
	DriverKafka = "kafka"

	// DriverRedis defines the driver type when integrating with a Redis queue.
	DriverRedis = "redis"
)

// Worker runtime drivers.
const (
	// DriverDocker defines the driver type when integrating with a Docker runtime.
	DriverDocker = "docker"

	// DriverKubernetes defines the driver type when integrating with a Kubernetes runtime.
	DriverKubernetes = "kubernetes"
)

// Server and worker secret drivers.
const (
	// DriverNative defines the driver type when integrating with a Vela secret service.
	DriverNative = "native"

	// DriverVault defines the driver type when integrating with a Vault secret service.
	DriverVault = "vault"
)

// Server source drivers.
const (
	// DriverGitHub defines the driver type when integrating with a Github source code system.
	DriverGithub = "github"

	// DriverGitLab defines the driver type when integrating with a Gitlab source code system.
	DriverGitlab = "gitlab"
)

// Server storage drivers.
const (
	// DriverMinio defines the driver type when integrating with a local storage system.
	DriverMinio = "minio"
)
