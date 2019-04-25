run:
	go run .

build:
	CGO_ENABLED=0 GOOS=linux go build -o app -a -ldflags '-extldflags "-static"' .

### Docker commands
docker-run:
	docker run --rm -p 8080:8080 stefanoschrs/fb-picture-cacher

docker-build:
	docker build -t stefanoschrs/fb-picture-cacher .
