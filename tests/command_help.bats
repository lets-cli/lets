load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_help
}

@test "command_help: help contains description and options" {
    run lets help test
    assert_success
    assert_output --partial "Run tests"
}

@test "command_help: must add new line between description and options" {
    run lets help test2
    assert_success
    assert_output --partial "Run tests"
}
