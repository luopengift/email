MAIN = *.go
APP = $(shell basename `pwd`) #当前目录名称
VERSION	= 1.0.0

TIME = $(shell date "+%F %T")
GIT = $(shell git rev-parse HEAD)
PKG = github.com/luopengift/version

FLAG = "-X '${PKG}.VERSION=${VERSION}' -X '${PKG}.APP=${APP}' -X '${PKG}.TIME=${TIME}' -X '${PKG}.GIT=${GIT}'"
build: 
	go build -ldflags $(FLAG) -o ${APP} ${MAIN}