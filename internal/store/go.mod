module github.com/rafaelcalleja/go-kit/internal/store

go 1.17

replace (
	github.com/rafaelcalleja/go-kit/internal/common => ./../common
	github.com/rafaelcalleja/go-kit/logger => ./../../logger
	github.com/rafaelcalleja/go-kit/uuid => ./../../uuid
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/golang/protobuf v1.5.2
	github.com/huandu/go-sqlbuilder v1.13.0
	github.com/rafaelcalleja/go-kit/internal/common v0.0.0-00010101000000-000000000000
	github.com/rafaelcalleja/go-kit/uuid v0.0.0-20220114085949-e6ff973b8411
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.43.0
)

require (
	cloud.google.com/go v0.34.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
