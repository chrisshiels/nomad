# Overview

Hashicorp Nomad-based container orchestration with Docker Engine,
Docker Registry, Dnsmasq, Consul, Nomad and Fabio.
```
        vm1:        vm2:        vm3:

        Fabio       Fabio       Fabio

        Nomad       Nomad       Nomad

        Consul      Consul      Consul

        Dnsmasq     Dnsmasq     Dnsmasq

        Docker
        Registry

        Docker      Docker      Docker
        Engine      Engine      Engine
```


# Downloads

```
host$ ( cd downloads ; make )
```


# Virtual machines

```
host$ export VAGRANT_DEFAULT_PROVIDER=libvirt
host$ sudo iptables -I INPUT 1 -i virbr+ -j ACCEPT
host$ vagrant up
host$ vagrant status

host$ vagrant ssh vm1

host$ vagrant ssh vm2

host$ vagrant ssh vm3
```


# Docker Engine 17.06.0.ce

```
vm123$ sudo yum install -y yum-utils device-mapper-persistent-data lvm2
vm123$ sudo yum-config-manager \
        --add-repo https://download.docker.com/linux/centos/docker-ce.repo
vm123$ sudo yum -y install docker-ce-17.06.0.ce

vm123$ ip=$(/sbin/ip addr list eth1 | \
        awk '/inet / { gsub(/\/.*/, "", $2); print $2 }')

vm123$ sudo mkdir /etc/docker
vm123$ sudo tee /etc/docker/daemon.json <<eof
{
  "storage-driver": "devicemapper",
  "dns": [
    "$ip"
  ],
  "dns-search": [
    "fabio",
    "service.consul"
  ]
}
eof

vm123$ sudo /bin/systemctl enable docker
vm123$ sudo /bin/systemctl start docker
vm123$ sudo /bin/systemctl status docker
```


# Docker Compose

This looks to be a PyInstaller packaged Python application.
```
vm1$ sudo install --owner root --group root --mode 755 \
        /vagrant/downloads/docker-compose-Linux-x86_64 \
        /usr/local/bin/docker-compose
```


# Docker Registry

Certificate generation:
```
vm1$ ( cd /vagrant/ssl/ ; make )
```


Docker Engine configuration:
```
vm123$ sudo mkdir -p /etc/docker/certs.d/vm1:5000
vm123$ sudo cp /vagrant/ssl/vm1.crt /etc/docker/certs.d/vm1:5000/ca.crt
vm123$ sudo /bin/systemctl restart docker.service
vm123$ sudo /bin/systemctl status docker.service
```


Start Docker Registry:
```
vm1$ cd /vagrant/registry/
vm1$ cat docker-compose.yml
registry:
  restart: always
  image: registry:2
  ports:
    - 5000:5000
  environment:
    REGISTRY_HTTP_TLS_CERTIFICATE: /certs/vm1.crt
    REGISTRY_HTTP_TLS_KEY: /certs/vm1.key
  volumes:
    - /vagrant/ssl:/certs

vm1$ sudo /usr/local/bin/docker-compose up -d
vm1$ sudo /usr/local/bin/docker-compose ps
```


Test Docker Registry:
```
vm1$ sudo docker pull alpine:latest
vm1$ sudo docker tag alpine:latest vm1:5000/alpine:latest
vm1$ sudo docker push vm1:5000/alpine:latest
```


# Dnsmasq

```
vm123$ sudo yum -y install dnsmasq
vm123$ sudo tee /etc/dnsmasq.d/custom.conf <<eof
address=/fabio/$ip
server=/consul/127.0.0.1#8600
eof

vm123$ sudo /bin/systemctl enable dnsmasq.service
vm123$ sudo /bin/systemctl start dnsmasq.service
vm123$ sudo /bin/systemctl status dnsmasq.service
```

Note:  Could enhance this further with Consul-template generating per-service
configuration for, e.g. date.fabio, time.fabio and web.fabio.


# Consul 0.8.4

```
vm123$ sudo yum -y install unzip

vm123$ sudo unzip -d /usr/local/bin/ \
        /vagrant/downloads/consul_0.8.4_linux_amd64.zip

vm123$ ip=$(/sbin/ip addr list eth1 | \
        awk '/inet / { gsub(/\/.*/, "", $2); print $2 }')

vm123$ sudo mkdir /etc/consul
vm123$ sudo tee /etc/consul/config.json <<eof
{
  "server": true,
  "bind_addr": "0.0.0.0",
  "client_addr": "0.0.0.0",
  "advertise_addr": "$ip",
  "bootstrap_expect": 3,
  "retry_join": [ "192.168.40.11", "192.168.40.12", "192.168.40.13" ],
  "data_dir": "/var/lib/consul"
}
eof

vm123$ sudo sh -c '/usr/local/bin/consul agent -config-dir /etc/consul/ >> /var/log/consul.log 2>&1' &

vm123$ consul members
```


# Nomad 0.5.6

```
vm123$ sudo unzip -d /usr/local/bin/ \
        /vagrant/downloads/nomad_0.5.6_linux_amd64.zip

vm123$ sudo mkdir /etc/nomad
vm123$ sudo tee /etc/nomad/config.hcl <<eof
log_level = "DEBUG"

datacenter = "dc1"

data_dir = "/var/lib/nomad"

server {
    enabled = true
    bootstrap_expect = 3
}

client {
    enabled = true
}
eof

vm123$ sudo sh -c '/usr/local/bin/nomad agent -config /etc/nomad/ >> /var/log/nomad.log 2>&1' &

vm123$ nomad server-members
```


# Fabio 1.5.0

```
vm123$ sudo install --owner root --group root --mode 755 \
        /vagrant/downloads/fabio-1.5.0-go1.8.3-linux_amd64 \
        /usr/local/bin/fabio

vm123$ sudo mkdir /etc/fabio
vm123$ sudo tee /etc/fabio/fabio.properties <<eof
log.access.target = stdout

ui.addr = :9998

proxy.addr = :80
eof

vm123$ sudo sh -c '/usr/local/bin/fabio -cfg /etc/fabio/fabio.properties >> /var/log/fabio-access.log 2> /var/log/fabio.log' &
```

To view Web user interface:
```
>> http://vm1:9998/
```


# Build images

```
vm1$ sudo yum -y install golang

vm1$ cd /vagrant/images/date
vm1$ make VERSION=1.0.0 build
vm1$ sudo docker tag cs/date:1.0.0 vm1:5000/cs/date:1.0.0
vm1$ sudo docker push vm1:5000/cs/date:1.0.0

vm1$ cd /vagrant/images/time
vm1$ make VERSION=1.0.0 build
vm1$ sudo docker tag cs/time:1.0.0 vm1:5000/cs/time:1.0.0
vm1$ sudo docker push vm1:5000/cs/time:1.0.0

vm1$ cd /vagrant/images/web
vm1$ make VERSION=1.0.0 build
vm1$ sudo docker tag cs/web:1.0.0 vm1:5000/cs/web:1.0.0
vm1$ sudo docker push vm1:5000/cs/web:1.0.0
```


# Run instances

```
vm1$ cd /vagrant/nomad

vm1$ cat date.nomad
vm1$ nomad run ./date.nomad

vm1$ cat time.nomad
vm1$ nomad run ./time.nomad

vm1$ cat web.nomad
vm1$ nomad run ./web.nomad


vm1$ nomad status
vm1$ nomad status date
vm1$ nomad status time
vm1$ nomad status web
```


# Checks

```
vm1$ curl -s http://vm1:8500/v1/catalog/services | python -mjson.tool
vm1$ curl -s http://vm1:8500/v1/catalog/service/date | python -mjson.tool
vm1$ curl -s http://vm1:8500/v1/catalog/service/time | python -mjson.tool
vm1$ curl -s http://vm1:8500/v1/catalog/service/web | python -mjson.tool

vm1$ sudo yum -y install bind-utils
vm1$ dig @vm1 date.service.consul in any
vm1$ dig @vm1 date.service.consul in srv
vm1$ dig @vm1 time.service.consul in any
vm1$ dig @vm1 time.service.consul in srv
vm1$ dig @vm1 web.service.consul in any
vm1$ dig @vm1 web.service.consul in srv

vm1$ curl --header 'Host: date' http://localhost/date ; echo
{"date":"20170622","hostname":"a66af5f2ae6d","version":"1.0.0"}

vm1$ curl --header 'Host: time' http://localhost/time ; echo
{"time":"16:29:37","hostname":"fdfdbc6e7720","version":"1.0.0"}

vm1$ curl --header 'Host: web' http://localhost/ ; echo
cfcebf69fe54 1.0.0:
20170622 - 0b7ce2689b1a 1.0.0
16:29:57 - 7002ddd70e12 1.0.0


vm1$ sudo docker ps -a | grep web
ef8a1a57532d        e85699e1b8e9        "/bin/sh -c '/web ..."   21 minutes ago      Up 21 minutes       192.168.121.111:31843->7000/tcp, 192.168.121.111:31843->7000/udp   web-2b7d2f1f-98b0-0224-8ea9-3a60110d4d2e

vm1$ sudo docker exec -i -t ef8a1a57532d /bin/sh

/ # cat /etc/resolv.conf
search fabio service.consul
nameserver 192.168.40.11

/ # getent hosts date.fabio
192.168.40.11     date.fabio  date.fabio

/ # getent hosts date.service.consul
192.168.121.102   date.service.consul  date.service.consul

/ # apk update
/ # apk add curl

/ # curl http://date/date ; echo
{"date":"20170628","hostname":"de74c9059586","version":"1.0.0"}

/ # curl http://time/time ; echo
{"time":"23:33:16","hostname":"3a0299c27ad3","version":"1.0.0"}

/ # curl http://web/ ; echo
94a5c24bbc4f 1.0.0:
20170628 - 9f656c81cc16 1.0.0
23:33:22 - 6bc21a4f52ab 1.0.0
```
Note Docker container instances are named *imagename-allocationid*.
