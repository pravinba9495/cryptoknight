eth:
	npm install
	npm run pkg
build: eth
	go build -o bin/kryptonite main.go
docker:
	npm install
	npm run docker
	GOOS=linux GOARCH=amd64 go build -o bin/kryptonite main.go
	docker build . -t pravinba9495/kryptonite:latest