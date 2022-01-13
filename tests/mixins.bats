setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/mixins
}

@test "mixins: mixins works" {
    run lets hello-from-minix
    assert_success
    assert_line --index 0 "Hello"
}
