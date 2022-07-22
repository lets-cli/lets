load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/no_depends
}

@test "no_depends: should skip depends for running command" {
    run lets ping-pong
    assert_success
    assert_line --index 0 "Ping"
    assert_line --index 1 "Pong"
    assert_line --index 2 "Done"

    run lets --no-depends ping-pong
    assert_success
    assert_line --index 0 "Done"
}
