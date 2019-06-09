LDFLAGS=-ldflags "-X main.version=`git rev-list -1 HEAD`"

default:
	go generate
	go build ${LDF
