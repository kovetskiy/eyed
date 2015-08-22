#!/bin/bash

tests_ensure eyed_run ":8080" "$TEST_DIR/reports/"
eyed_task_id=$(cat `tests_stdout`)
eyed_stderr=$(tests_background_stderr $eyed_task_id)

# directory should be created automatically
tests_test -d "$TEST_DIR/reports/"

tests_ensure curl -vs localhost:8080/
tests_assert_stderr_re '400 Bad Request'

# 1 for . directory
tests_test $(find "$TEST_DIR/reports" | wc -l) -eq 1

tests_ensure curl -vs localhost:8080/foo/bar
tests_assert_stderr_re '400 Bad Request'

# 1 for . directory
tests_test $(find "$TEST_DIR/reports" | wc -l) -eq 1
