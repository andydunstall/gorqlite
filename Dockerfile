FROM golang:1.17

RUN wget -O toxiproxy-2.1.4.deb https://github.com/Shopify/toxiproxy/releases/download/v2.1.4/toxiproxy_2.1.4_amd64.deb
RUN dpkg -i toxiproxy-2.1.4.deb

WORKDIR /usr/local
RUN wget https://github.com/rqlite/rqlite/releases/download/v6.7.0/rqlite-v6.7.0-linux-amd64.tar.gz && \
	tar -zxf rqlite-v6.7.0-linux-amd64.tar.gz && \
	cp rqlite-v6.7.0-linux-amd64/* /usr/local/bin
