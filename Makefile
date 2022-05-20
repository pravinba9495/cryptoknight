build:
	npm install
	npm run tsc
docker:
	make build
	docker build . -t pravinba9495/kryptonite:latest