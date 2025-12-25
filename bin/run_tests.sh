#!/bin/bash

go build ./iteng
#go test ./iteng -v

cd iteng
go test -v

