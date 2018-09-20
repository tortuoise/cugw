#!/bin/sh
protoc --go_out=plugins=grpc:. cugw.proto


