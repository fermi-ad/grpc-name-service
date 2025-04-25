#!/bin/bash

TOP=../../
protodir=$TOP/nameserver/src/main/proto
targetdir=.
PROTOFILES=`ls $protodir/*.proto`
protoc -I=${protodir} --go_out=${targetdir} --go-grpc_out=${targetdir} $PROTOFILES
