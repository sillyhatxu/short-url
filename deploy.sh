go clean
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
docker build -t short-url .
docker tag short-url:latest xushikuan/short-url:1.0
docker push xushikuan/short-url:1.0
