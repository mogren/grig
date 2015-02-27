build: grig

grig: *.go
	go build -o $@ $^
