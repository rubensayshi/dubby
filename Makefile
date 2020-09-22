#
# Makefile is purely for DEVELOPMENT purposes
#
RUN_=
ifdef RUN
	RUN_=-run $(RUN)
endif

PARALLEL_=-parallel=8
ifdef PARALLEL
	PARALLEL_=-parallel $(PARALLEL)
endif

COVERPROFILE_=
ifdef COVERPROFILE
	COVERPROFILE_=-coverprofile=$(COVERPROFILE) -coverpkg=./src/...
endif

TESTS ?= ./src/...
EXTRA ?=

#
# reflex-test
# this runs `go test` using reflex to auto reload
#
# there's a few possible arguments to control `go tests` that are passed down
# 	TESTS=		can be used to specify which tests to run as a path, defaults to ./src/...
# 				you can do `TESTS=./src/modules/orm/...`
# 				or even multiple (note the '') `TESTS='./src/modules/orm/... ./src/container/...'`
#
#	PARALLEL=	sets the -parellel=<X> number, defaults to 8
#	RUN=		sets the -run flag to apply a regex to which tests to run (eg; `RUN=Test_apiLogin1`)
#	EXTRA=		everything in EXTRA is passed down to go test (eg; `EXTRA=-v` or `EXTRA=`-v -race -cover`)
#
reflex-test:
	reflex -s -r '^src/.*\.go$$' -R 'vendor/' -R 'tmp/' -- sh -c 'make go-test'

#
# go-test
# this simply runs `go test` without reflex
# same args as reflex-test can be used
#
go-test:
	go test -failfast $(PARALLEL_) $(EXTRA) $(TESTS) $(RUN_) $(COVERPROFILE_) && echo "TESTS DONE";

BINDIR ?= ./bin

build: build-dujsonimporter build-dujsonexporter

build-dujsonimporter:
	go build -o $(BINDIR)/dujsonimporter src/cmd/dujsonimporter/main.go

build-dujsonexporter:
	go build -o $(BINDIR)/dujsonexporter src/cmd/dujsonexporter/main.go

build-release:
	rm -rf ./release
	mkdir ./release
	GOOS=windows GOARCH=amd64 go build -o release/dujsonimporter.exe src/cmd/dujsonimporter/main.go
	GOOS=windows GOARCH=amd64 go build -o release/dujsonexporter.exe src/cmd/dujsonexporter/main.go
	GOOS=darwin GOARCH=amd64 go build -o release/dujsonimporter-darwin src/cmd/dujsonimporter/main.go
	GOOS=darwin GOARCH=amd64 go build -o release/dujsonexporter-darwin src/cmd/dujsonexporter/main.go
	GOOS=linux GOARCH=amd64 go build -o release/dujsonimporter-linux src/cmd/dujsonimporter/main.go
	GOOS=linux GOARCH=amd64 go build -o release/dujsonexporter-linux src/cmd/dujsonexporter/main.go
