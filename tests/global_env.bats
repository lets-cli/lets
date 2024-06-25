load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/global_env
}

@test "global_env: should provide env to command" {
    run lets global-env
    assert_success
    assert_line --index 0 "INT=1"
    assert_line --index 1 "STR=hi"
    assert_line --index 2 "STR_INT=1"
    assert_line --index 3 "BOOL=true"
    assert_line --index 4 "ORIGINAL=b"
    assert_line --index 5 "BAR=Bar"
    assert_line --index 6 "FOO=bb1da47569d9fbe3b5f2216fdbd4c9b040ccb5c1"
}
