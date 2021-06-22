#!/usr/bin/env bash

tag=$(git describe --tags)
builds=(
  darwin/amd64
  darwin/arm64
  freebsd/386
  freebsd/amd64
  freebsd/arm
  freebsd/arm64
  linux/386
  linux/amd64
  linux/arm
  linux/arm64
  windows/386
  windows/amd64
  windows/arm
)

for pair in ${builds[@]}; do
  os=`dirname $pair`
  arch=`basename $pair`
  dir=dist/titania.go-$os-$arch-$tag
  echo creating: $dir
  if [ $os = windows ]; then
    GOOS=$os GOARCH=$arch go build -o $dir/titania.go.exe cmd/titania.go/*
    zip -j $dir.zip $dir/titania.go.exe
  else
    GOOS=$os GOARCH=$arch go build -o $dir/titania.go     cmd/titania.go/*
    zip -j $dir.zip $dir/titania.go
  fi
  rm -rf $dir
done
