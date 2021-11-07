FROM kong/go-plugin-tool:latest-alpine-latest AS builder

RUN mkdir -p /tmp/key-checker/

COPY . /tmp/key-checker/

RUN cd /tmp/key-checker/ && \
    go get github.com/Kong/go-pdk && \
    go mod init kong-go-plugin && \
    go get -d -v github.com/Kong/go-pluginserver && \
    go build github.com/Kong/go-pluginserver && \
    go get -d -v github.com/basgys/goxml2json && \
    go build github.com/basgys/goxml2json && \
    go build -buildmode plugin key-checker.go

FROM kong:2.6.0-alpine

RUN mkdir /tmp/go-plugins
COPY --from=builder  /tmp/key-checker/go-pluginserver /usr/local/bin/go-pluginserver
COPY --from=builder  /tmp/key-checker/key-checker.so /tmp/go-plugins
COPY config.yml /tmp/config.yml

USER root
RUN chmod -R 777 /tmp
#RUN chmod -R 777 /usr/local/bin
#RUN chmod -R 777 /usr/local/kong
RUN /usr/local/bin/go-pluginserver -version && \
    cd /tmp/go-plugins && \
    /usr/local/bin/go-pluginserver -dump-plugin-info key-checker
USER kong