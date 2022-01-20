build:
	npm install
	npm run tsc
	npm run pkg
docker:
	npm install
	npm run tsc
	npm run docker
	docker build . -t pravinba9495/kryptonite:latest