#!/bin/bash

function usage() {
    echo "$0 <protopath>"
}

TOP=..    
protodir=$TOP/nameserver/src/main/proto
targetdir=.
PROTOFILES=`ls $protodir/*.proto`
protoc -I=${protodir} --go_out=${targetdir} --go-grpc_out=${targetdir} $PROTOFILES
