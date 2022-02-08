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
    assert_line --index 2 "BAR=Bar"
    assert_line --index 3 "FOO=bb1da47569d9fbe3b5f2216fdbd4c9b040ccb5c1"
}
