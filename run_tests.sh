#!/bin/bash

EYED_BIN="$(readlink -f eyed-test)"
go build -o "$EYED_BIN"
if [ $? -ne 0 ]; then
    exit 1
fi

# bash tests library
if [ ! -f tests/lib/tests.sh ]; then
    git submodule init
    git submodule update

    if [ ! -f tests/lib/tests.sh ]; then
        echo "file 'tests/lib/tests.sh' not found"
        exit 1
    fi
fi

source tests/lib/tests.sh
source tests/functions.sh

cd tests
TEST_VERBOSE=10
tests_run_all
