#!/bin/bash

protodir=.
targetdir=./target

PROTOFILES=`ls ${protodir}/*.proto`

set -x
protoc -I=${protodir} --go_out=${targetdir} --go-grpc_out=${targetdir} ${PROTOFILES}
