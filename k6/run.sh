#!/bin/bash

FILEPATH=$(readlink -f "$0")
BASEPATH=$(dirname "$FILEPATH")

k6 run $BASEPATH/*.js -v
