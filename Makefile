all: build

.PHONY: build

build:
	cd ./gate && go build
	cd ./notifer && go build
	cd ./rest_api && go build
