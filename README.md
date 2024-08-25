# go playground

A place I try out random go stuff.

### Albums web applicaiton

#### Build

```shell
go build ./cmd/ws/album.go 
```

#### Run

```shell
docker build -t ws-albums -f Dockerfile.ws . && docker run --rm -d -p 8080:8080 ws-albums
```
