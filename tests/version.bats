setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
}

@test "version: show lets version for -v" {
    run lets -v
    assert_success
    assert_line --index 0 "lets version 0.0.0-dev"
}

@test "version: show lets version for --version" {
    run lets --version
    assert_success
    assert_line --index 0 "lets version 0.0.0-dev"
}
