build: grig

install:
	- go get github.com/mogren/grig/vose

grig: *.go
	go build -o $@ $^
