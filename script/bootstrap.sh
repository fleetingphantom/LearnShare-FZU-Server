#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=LearnShare
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}
