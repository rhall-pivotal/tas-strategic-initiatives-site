#!/bin/bash

URL=$1
LOCAL_FILE=$2

aria2c -x 5 --out=$LOCAL_FILE $URL
