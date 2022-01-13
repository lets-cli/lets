load test_helpers

setup() {
    load "${BATS_UTILS_PATH}/bats-support/load.bash"
    load "${BATS_UTILS_PATH}/bats-assert/load.bash"
    cd ./tests/global_before
}

@test "global_before: should insert before script for each cmd" {
    run lets hello
    assert_success
    assert_line --index 0 "Hello"
}

@test "global_before: should merge before scripts from mixins" {
    run lets world
    assert_success
    assert_line --index 0 "World"
}
