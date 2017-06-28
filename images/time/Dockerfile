# 'Dockerfile'.
# Chris Shiels.


from alpine:latest

expose 7002

env PORT=7002

run mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

cmd /time -p $PORT

copy time /
