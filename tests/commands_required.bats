setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/commands_required
}

@test "commands_required: fail if no commands in lets config" {
    run lets
    assert_failure
    assert_line --index 0 "lets: config error: 'commands' can not be empty"
}
