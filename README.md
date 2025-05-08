# delve

The smallest AI wrapper?

## Build from scratch

```
docker build --tag=delve --file=build.Dockerfile .
docker create --name delve delve
docker cp delve:/app/delve delve
docker rm delve
```
