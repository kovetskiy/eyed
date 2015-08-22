#!/bin/bash

tests_ensure eyed_run ":54321" "$TEST_DIR/reports/"
eyed_task_id=$(cat `tests_stdout`)
eyed_stderr=$(tests_background_stderr $eyed_task_id)

# directory should be created automatically
tests_test -d "$TEST_DIR/reports/"

tests_ensure curl -vs localhost:54321/foo
tests_assert_stderr_re '200 OK'

tests_test -f "$TEST_DIR/reports/foo"

tests_do rm "$TEST_DIR/reports/foo"
tests_do mkdir "$TEST_DIR/reports/foo"

tests_ensure curl -vs localhost:54321/foo
tests_assert_stderr_re '500 Internal Server Error'

# error should be displayed as http response
tests_assert_stdout_re 'is a directory'
