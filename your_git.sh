#!/bin/sh
#
# DON'T EDIT THIS!
#
# CodeCrafters uses this file to test your code. Don't make any changes here!
#
# DON'T EDIT THIS!
set -e
tmpFile=$(mktemp)

cwd=$PWD

cd $(dirname $0)/app && go build -o "$tmpFile" && cd $cwd

exec "$tmpFile" "$@"
