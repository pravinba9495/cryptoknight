build:
	npm install
	npm run tsc
docker:
	npm install
	npm run tsc
	docker build . -t pravinba9495/kryptonite:latest