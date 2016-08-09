#!/bin/sh
#
# build.sh
# Build a project
# By J. Stuart McMurray
# Created 20160221
# Last Modified 20160808

set -e

PROG=$(basename $(pwd))

echo "Building $PROG"

go vet

for GOOS in linux openbsd darwin; do
        for GOARCH in 386 amd64; do
                export GOOS GOARCH
                N="$PROG.$GOOS.$GOARCH"
                go build -o "$N"
                ls -l $N
        done
done

echo Done.
