module github.com/rafaelcalleja/go-kit/cmd/cobra/version

go 1.17

replace (
	github.com/rafaelcalleja/go-kit/cmd/helper => ../../../cmd/helper
	github.com/rafaelcalleja/go-kit/cmd/termcolor => ../../termcolor
	github.com/rafaelcalleja/go-kit/logger => ../../../logger
)

require (
	github.com/rafaelcalleja/go-kit/cmd/helper v0.0.0-00010101000000-000000000000
	github.com/rafaelcalleja/go-kit/cmd/termcolor v0.0.0-00010101000000-000000000000
	github.com/rafaelcalleja/go-kit/logger v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.2.0
)

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jenkins-x/jx-helpers/v3 v3.1.3 // indirect
	github.com/jenkins-x/jx-logging/v3 v3.0.6 // indirect
	github.com/jenkins-x/logrus-stackdriver-formatter v0.2.3 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rickar/props v0.0.0-20170718221555-0b06aeb2f037 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
)
