module github.com/rafaelcalleja/go-kit/internal/store

go 1.17

replace (
	github.com/rafaelcalleja/go-kit/internal/common => ./../common
	github.com/rafaelcalleja/go-kit/logger => ./../../logger
	github.com/rafaelcalleja/go-kit/uuid => ./../../uuid
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/deepmap/oapi-codegen v1.9.1
	github.com/go-chi/chi/v5 v5.0.7
	github.com/golang/protobuf v1.5.2
	github.com/huandu/go-sqlbuilder v1.13.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rafaelcalleja/go-kit/internal/common v0.0.0-00010101000000-000000000000
	github.com/rafaelcalleja/go-kit/uuid v0.0.0-20220114085949-e6ff973b8411
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.43.0
)

require (
	cloud.google.com/go/compute v1.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-chi/cors v1.2.0 // indirect
	github.com/go-chi/render v1.0.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/net v0.0.0-20210913180222-943fd674d43e // indirect
	golang.org/x/sys v0.0.0-20220204135822-1c1b9b1eba6a // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220204002441-d6cc3cc0770e // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/gorm v1.22.5 // indirect
)
