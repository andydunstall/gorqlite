# Docker environment to run system tests (`make system-test`).
FROM golang:1.17

# Install rqlite.
WORKDIR /usr/local
RUN wget https://github.com/rqlite/rqlite/releases/download/v6.7.0/rqlite-v6.7.0-linux-amd64.tar.gz && \
	tar -zxf rqlite-v6.7.0-linux-amd64.tar.gz && \
	cp rqlite-v6.7.0-linux-amd64/* /usr/local/bin
