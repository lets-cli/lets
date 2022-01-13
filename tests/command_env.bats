load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_env
}

@test "command_env: should provide env to command" {
    run lets env
    assert_success
    assert_line --index 0 "ONE=1"
    assert_line --index 1 "TWO=two"
}
