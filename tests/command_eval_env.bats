load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_eval_env
}

@test "command_eval_env: should compute and provide env to command" {
    run lets eval-env
    assert_success
    assert_line --index 0 "ONE=1"
    assert_line --index 1 "TWO=two"
    assert_line --index 2 "COMPUTED=Computed env"
}
