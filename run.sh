#!/usr/bin/env bash

TOPIC=$1
CONSUMER_GROUP=$2

ECHO "[topic-bleed] compiling..."
go build -o lib/topic-bleed
ECHO "[topic-bleed] complete"

ECHO "[topic-bleed] starting: topic=$1 consumer-group=$2"
./lib/topic-bleed -topic=$1 -group=$2