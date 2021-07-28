module github.com/go-vela/server

go 1.15

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/aws/aws-sdk-go v1.38.57
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/color v1.10.0 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-vela/compiler v0.8.2-0.20210728151243-38af4aebce57
	github.com/go-vela/pkg-queue v0.8.1
	github.com/go-vela/types v0.8.3-0.20210726122150-0eaf6091307b
	github.com/google/go-cmp v0.5.6
	github.com/google/go-github/v37 v37.0.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/hashicorp/vault/api v1.1.1
	github.com/jackc/pgproto3/v2 v2.0.7 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20210503195802-e9a32991a82e // indirect
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637
	gorm.io/driver/postgres v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
)
