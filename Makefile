all: build

.PHONY: build

build:
	cd ./data && go build
	cd ./gate && go build
	cd ./notifer && go build
	cd ./rest_api && go build
	cd ./session && go build
