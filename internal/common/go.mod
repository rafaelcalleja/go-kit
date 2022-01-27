module github.com/rafaelcalleja/go-kit/internal/common

go 1.17

replace github.com/rafaelcalleja/go-kit/uuid => ./../../uuid

require (
	cloud.google.com/go v0.34.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/petermattis/goid v0.0.0-20220111183729-e033e1e0bdb5
	github.com/pkg/errors v0.8.1
	github.com/rafaelcalleja/go-kit/uuid v0.0.0-20220114085949-e6ff973b8411
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
