# SuperCalculator

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/wajox/gobase)](https://goreportcard.com/report/github.com/wajox/gobase)
[![codecov](https://codecov.io/gh/wajox/gobase/branch/master/graph/badge.svg?token=0K79C2LH2K)](https://codecov.io/gh/wajox/gobase)
[![Build Status](https://travis-ci.org/wajox/gobase.svg?branch=master)](https://travis-ci.org/wajox/gobase)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

This is a simple skeleton for golang application. Inspired by development experience and updated according to github.com/golang-standards/project-layout.

## How to use?



## Structure

* /front-end - 
* /back-end -
  * /agent
  * /orkestrator

## Commands
```sh
# install dev tools(wire, golangci-lint, swag, ginkgo)
make install-tools

# start test environment from docker-compose-test.yml
make start-docker-compose-test

# stop test environment from docker-compose-test.yml
make stop-docker-compose-test

# run all tests
make test-all

# run go generate
make gen

# generate source code from .proto files
make proto

# generate dependencies with wire
make deps
```


## Tools and packages
* gin-gonic
* ginkgo with gomega
* spf13/viper
* spf13/cobra
* envy
* zerolog
* golangci-lint
* wire
* swag
* migrate
* protoc
* jsonapi
* docker with docker-compose
