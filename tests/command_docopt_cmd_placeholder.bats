load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_docopt_cmd_placeholder
}

@test "command_docopt_cmd_placeholder: should run with docopt from yaml alias" {
    # We can use yaml alias syntax to prevent repetition of docsopt description
    # The placeholder string is \$\{LETS_COMMAND_NAME\} for now
    run lets cmd-1 posarg --config=some_path
    assert_success
    assert_line --index 0 "posarg"
    assert_line --index 1 "some_path"
}

@test "command_docopt_cmd_placeholder: should fail with docopt from yaml alias wo placeholder" {
    run lets cmd-2 posarg --config=some_path

    assert_failure
    assert_line --index 0 --partial "no such option"
}
