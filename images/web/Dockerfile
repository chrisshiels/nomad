# 'Dockerfile'.
# Chris Shiels.


from alpine:latest

expose 7000

env PORT=7000
env DATEENDPOINT=date.service.consul
env TIMEENDPOINT=time.service.consul

run mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

cmd /web -p $PORT -dateendpoint $DATEENDPOINT -timeendpoint $TIMEENDPOINT

copy web /
