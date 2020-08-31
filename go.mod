module github.com/go-vela/server

go 1.13

replace github.com/go-vela/compiler => github.com/JordanSussman/compiler v0.1.3-0.20200818133900-c4a5e4bebff2

replace github.com/go-vela/types => github.com/JordanSussman/types v0.1.2-0.20200817233341-2853ce2956de

require (
	github.com/coreos/go-semver v0.3.0
	github.com/denisenkom/go-mssqldb v0.0.0-20191128021309-1d7a30a10f73 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/frankban/quicktest v1.7.2 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-vela/compiler v0.5.1
	github.com/go-vela/types v0.5.1
	github.com/google/go-github/v29 v29.0.3
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-hclog v0.10.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.6 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/vault/api v1.0.4
	github.com/jinzhu/gorm v1.9.14
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.6.0
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637
)
