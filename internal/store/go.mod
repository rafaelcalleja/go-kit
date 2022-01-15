module github.com/rafaelcalleja/go-kit/internal/store

go 1.17

replace (
	github.com/rafaelcalleja/go-kit/internal/common => ./../common
	github.com/rafaelcalleja/go-kit/logger => ./../../logger
	github.com/rafaelcalleja/go-kit/uuid => ./../../uuid
)

require (
	github.com/golang/protobuf v1.5.2
	github.com/rafaelcalleja/go-kit/internal/common v0.0.0-00010101000000-000000000000
	github.com/rafaelcalleja/go-kit/logger v0.0.0-20220106180013-2a82d5d5e135
	github.com/rafaelcalleja/go-kit/uuid v0.0.0-20220114085949-e6ff973b8411
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.43.0
)

require (
	cloud.google.com/go v0.34.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/jenkins-x/jx-logging/v3 v3.0.6 // indirect
	github.com/jenkins-x/logrus-stackdriver-formatter v0.2.3 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rickar/props v0.0.0-20170718221555-0b06aeb2f037 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
