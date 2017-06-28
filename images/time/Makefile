# 'Makefile'.
# Chris Shiels.


VERSION = latest


build:	Dockerfile time.go
	go build -ldflags "-X main.version=$(VERSION)" time.go
	sudo docker build -t cs/time:$(VERSION) .
	touch $@


run:	build
	sudo docker run --rm -i -t -p 7001:7001 cs/time:$(VERSION)
