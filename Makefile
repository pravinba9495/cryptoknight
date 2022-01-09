build:
	go build -o kryptonite main.go
docker:
	GOOS=linux GOARCH=amd64 go build -o kryptonite main.go
	docker build . -t pravinba9495/kryptonite:latest