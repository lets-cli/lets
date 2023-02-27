load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/global_eval_env
}

@test "global_eval_env: should compute env from eval_env and provide env to command" {
    run lets global-eval_env
    assert_success
    assert_line --index 0 "ONE=1"
    assert_line --index 1 "TWO=2"
}
