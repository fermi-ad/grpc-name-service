#!/bin/bash

set -x

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd $SCRIPT_DIR/../nsclient

function create() {
    ./nsclient create $1 --json="$2" 
}

function create_channels() {
    sysname=$1
    device=$2
    jsonend=',"accesscontrols":[{"role":"user", "read":true, "write":false}, {"role":"admin", "read":true, "write":true}]}'
    for c in `seq 1 5`; do
        json="{\"name\": \"ch${sysname}$c\", \"deviceName\":\"$device\""$jsonend
        echo $json
        ./nsclient create channel --json "$json"
    done
}
create location_type '{"name": "building"}'
create location_type '{"name": "room" }'
create location_type '{"name": "cabinet"}'
create location_type '{"name": "rack"}'

create alarm_type '{"name": "minor"}'
create alarm_type '{"name": "major"}'

create role '{"name": "user"}'
create role '{"name": "admin"}'

create location '{"name": "1234", "locationTypeName": "rack"}'
create location '{"name": "1235", "locationTypeName": "rack"}'

create node '{"hostname": "waf01", "ipAddress":"123.0.0.1", "locationName": "1234"}'
create node '{"hostname": "waf02", "ipAddress":"123.0.0.2", "locationName": "1234"}'
create node '{"hostname": "waf03", "ipAddress":"123.0.0.3", "locationName": "1234"}'

create device '{"name": "ioc01", "nodeHostname": "waf01"}'
create device '{"name": "ioc02", "nodeHostname": "waf01"}'
create device '{"name": "ioc03", "nodeHostname": "waf01"}'
create device '{"name": "ioc04", "nodeHostname": "waf01"}'
create device '{"name": "ioc05", "nodeHostname": "waf01"}'
create device '{"name": "ioc11", "nodeHostname": "waf02"}'
create device '{"name": "ioc12", "nodeHostname": "waf02"}'
create device '{"name": "ioc13", "nodeHostname": "waf02"}'
create device '{"name": "ioc14", "nodeHostname": "waf02"}'
create device '{"name": "ioc15", "nodeHostname": "waf02"}'

create_channels a ioc01
create_channels b ioc02
create_channels c ioc03
create_channels d ioc04
create_channels e ioc05
create_channels aa ioc11
create_channels ab ioc12
create_channels ac ioc13
create_channels ad ioc14
create_channels ae ioc15

