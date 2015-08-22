#!/bin/bash

tests_ensure eyed_run ":54321" "$TEST_DIR/reports/"
eyed_task_id=$(cat `tests_stdout`)
eyed_stderr=$(tests_background_stderr $eyed_task_id)

# directory should be created automatically
tests_test -d "$TEST_DIR/reports/"

tests_ensure curl -vs localhost:54321/foo
tests_assert_stderr_re '200 OK'

tests_test -f "$TEST_DIR/reports/foo"

tests_assert_re $eyed_stderr 'report.*foo.*hostname'

tests_ensure curl -vs localhost:54321/bar
tests_assert_stderr_re '200 OK'

tests_test -f "$TEST_DIR/reports/foo"
tests_test -f "$TEST_DIR/reports/bar"

tests_assert_re $eyed_stderr 'report.*foo.*hostname'
tests_assert_re $eyed_stderr 'report.*bar.*hostname'

tests_ensure curl -vs localhost:54321/foo
tests_assert_stderr_re '200 OK'

tests_test $(cat "$TEST_DIR/reports/foo" | wc -l) -eq 2
