#!/bin/sh

orgName="privatesquare"
imageName="bkst-users-api"
outputPath="docker-context"
version="1.0.0"

docker container rm -f ${imageName}

go fmt ./...
go get github.com/mitchellh/gox
env CGO_ENABLED=0 gox -os="linux" -arch="amd64" -output="${outputPath}/${imageName}-{{.OS}}-{{.Arch}}"

docker build -t ${orgName}/${imageName}:${version}  docker-context --build-arg VERSION=${version}
