# 'Makefile'.
# Chris Shiels.


TARGETS = consul_0.8.4_linux_amd64.zip \
	  docker-compose-Linux-x86_64 \
	  fabio-1.5.0-go1.8.3-linux_amd64 \
	  nomad_0.5.6_linux_amd64.zip


all: $(TARGETS)


consul_0.8.4_linux_amd64.zip:
	curl -L -O https://releases.hashicorp.com/consul/0.8.4/consul_0.8.4_linux_amd64.zip


docker-compose-Linux-x86_64:
	curl -L -O https://github.com/docker/compose/releases/download/1.14.0/docker-compose-Linux-x86_64


fabio-1.5.0-go1.8.3-linux_amd64:
	curl -L -O https://github.com/fabiolb/fabio/releases/download/v1.5.0/fabio-1.5.0-go1.8.3-linux_amd64


nomad_0.5.6_linux_amd64.zip:
	curl -L -O https://releases.hashicorp.com/nomad/0.5.6/nomad_0.5.6_linux_amd64.zip


clean:
	rm -vf $(TARGETS)
