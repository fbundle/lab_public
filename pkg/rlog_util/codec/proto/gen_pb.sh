#!/usr/bin/env bash
protoc --go_out=. --proto_path=./src ./src/rlog.proto