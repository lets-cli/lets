load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/command_after
}

@test "command_after: should run after script if cmd string" {
    run lets cmd-with-after
    assert_success
    assert_line --index 0 "Main"
    assert_line --index 1 "After"
}

@test "command_after: should run after script if cmd as map" {
    run lets cmd-as-map-with-after
    assert_success
    assert_line --index 0 "Main"
    assert_line --index 1 "After"
}

@test "command_after: should not shadow exit code from cmd" {
    run lets failure

    [[ $status = 113 ]]
    assert_line --index 0 "After"
}

@test "command_after: should not shadow exit code from cmd-as-map" {
    run lets failure-as-map

    [[ $status = 113 ]]
    assert_line --index 0 "After"
}
