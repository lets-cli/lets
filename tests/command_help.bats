load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_help
}


TEST_HELP_MESSAGE=$(cat <<EOF
Run tests
Unit tests are essention for success.

Example: lets test

Usage: lets test [<test_name>]
EOF
)

@test "command_help: help contains description and options" {
    run lets help test
    assert_success
    assert_output "${TEST_HELP_MESSAGE}"
}


TEST2_HELP_MESSAGE=$(cat <<EOF
Run tests

Usage: lets test2 [<test_name>]
EOF
)

@test "command_help: must add new line between description and options" {
    run lets help test2
    assert_success
    assert_output "${TEST2_HELP_MESSAGE}"
}