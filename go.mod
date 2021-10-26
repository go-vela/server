module github.com/go-vela/server

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/aws/aws-sdk-go v1.40.54
	github.com/fatih/color v1.10.0 // indirect
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-vela/compiler v0.10.1-0.20211025223007-bdcd7b7f8de0
	github.com/go-vela/pkg-queue v0.10.0
	github.com/go-vela/types v0.10.0
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/google/go-cmp v0.5.6
	github.com/google/go-github/v39 v39.2.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/hashicorp/vault/api v1.1.1
	github.com/joho/godotenv v1.4.0
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637
	gorm.io/driver/postgres v1.1.2
	gorm.io/driver/sqlite v1.1.6
	gorm.io/gorm v1.21.16
)
