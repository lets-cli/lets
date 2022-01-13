load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_shell
}

@test "command_shell: should run command using shell specified in command" {
    run lets show-shell
    assert_success
    assert_line --index 0 "/bin/sh"
}
