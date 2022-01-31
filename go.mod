module github.com/rafaelcalleja/go-kit

go 1.17

replace github.com/rafaelcalleja/go-kit/cmd/helper => ./cmd/helper

replace github.com/rafaelcalleja/go-kit/logger => ./logger

require (
	github.com/petermattis/goid v0.0.0-20220111183729-e033e1e0bdb5
	github.com/stretchr/testify v1.7.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
