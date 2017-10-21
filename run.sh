#!/usr/bin/env bash

ECHO "[topic-bleed] compiling..."
go build -o lib/topic-bleed

ECHO "[topic-bleed] starting:"
HUMAN_LOG=1 ./lib/topic-bleed