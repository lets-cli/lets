load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/global_eval_env
}

@test "global_eval_env: should fail because eval_env is not supported" {
    run lets global-eval_env
    assert_failure
    assert_output --partial "keyword 'eval_env' not supported"
}
