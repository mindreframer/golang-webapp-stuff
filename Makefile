#export GOPATH=/home/travis/gopath/jmadan/go-msgstory
#export GOBIN=/home/travis/gopath/bin

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w
GOGET=$(GOCMD) get

TARG=go-msgstory

GOFILES=\
	main.go\
	register/register.go\
	authenticate/authenticate.go\
	circle/circle.go\
	user/user.go\
	message/message.go\

all: test build

build:
	${GOBUILD} .
	

format:
	${GOFMT} -w ${GOFILES}

test:
	${GOTEST} ./user
	${GOTEST} ./message
	${GOTEST} ./register
	${GOTEST} ./authenticate
	${GOTEST} ./circle

.PHONY: advice build clean documentation format test