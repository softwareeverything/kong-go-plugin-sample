# kong-go-plugin-xml-to-json-transformator

This plugin basically transforms XML to JSON for Kong API Gateway.

```
docker build --no-cache -t kong-demo .
```

Images are built and we can run them as containers;

```
docker run -it --rm --name kong-go-plugins --network=kong-net --ip 172.18.0.20 \
  -e "KONG_GO_PLUGINS_DIR=/tmp/go-plugins" \
  -e "KONG_DECLARATIVE_CONFIG=/tmp/config.yml" \
  -e "KONG_PLUGINS=bundled,key-checker" \
  -e "KONG_DATABASE=postgres" \
  -e "KONG_PG_HOST=kong-database" \
  -e "KONG_PG_USER=kong" \
  -e "KONG_PG_PASSWORD=kong" \
  -e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
  -e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
  -e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
  -e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
  -e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
  -e "KONG_ADMIN_LISTEN=0.0.0.0:8001, 0.0.0.0:8444 ssl" \
  -e "KONG_LOG_LEVEL=debug" \
  -p 10000:8000 \
  -p 10443:8443 \
  -p 127.0.0.1:10001:8001 \
  -p 127.0.0.1:10444:8444 \
  kong-demo
```

Second option is;

```
docker run -it --rm --name kong-go-plugins \
  -e "KONG_DATABASE=off" \
  -e "KONG_GO_PLUGINS_DIR=/tmp/go-plugins" \
  -e "KONG_DECLARATIVE_CONFIG=/tmp/config.yml" \
  -e "KONG_PLUGINS=key-checker" \
  -e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
  -p 8000:8000 \
  kong-demo
```

(FROM kong:2.0.4-alpine)