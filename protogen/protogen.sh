#!/bin/sh

# Change to protofile directory
cd ../protofiles/

mkdir compiled

protoc --go_out=./compiled --plugin=../bin/protoc-gen-go transit.proto gtfs_realtime.proto nyct_realtime.proto