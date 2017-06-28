# 'Makefile'.
# Chris Shiels.


VERSION = latest


build:	Dockerfile date.go
	go build -ldflags "-X main.version=$(VERSION)" date.go
	sudo docker build -t cs/date:$(VERSION) .
	touch $@


run:	build
	sudo docker run --rm -i -t -p 7001:7001 cs/date:$(VERSION)
