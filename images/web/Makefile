# 'Makefile'.
# Chris Shiels.


VERSION = latest


build:	Dockerfile web.go
	go build -ldflags "-X main.version=$(VERSION)" web.go
	sudo docker build -t cs/web:$(VERSION) .
	touch $@


run:	build
	sudo docker run --rm -i -t -p 7001:7001 cs/web:$(VERSION)
