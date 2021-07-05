#!/bin/sh
#
# DON'T EDIT THIS!
#
# CodeCrafters uses this file to test your code. Don't make any changes here!
#
# DON'T EDIT THIS!
set -e
tmpFile=$(mktemp)

go env

cd app && go get && cd ..

go build -o "$tmpFile" $(dirname "$0")/app/*.go
exec "$tmpFile" "$@"
