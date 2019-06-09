LDFLAGS=-ldflags "-X main.version=`git describe`"

default:
	go generate
	GOOS=darwin GOARCH=amd64 go build -o goread-darwin-amd64 ${LDFLAGS}
	GOOS=linux GOARCH=amd64 go build -o goread-linux-amd64 ${LDFLAGS}
